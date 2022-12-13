package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rmarken5/cfcs/common"
	file_listener "github.com/rmarken5/cfcs/server/file-listener"
	"github.com/rmarken5/cfcs/server/observer"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:generate mockgen -destination=./mock_net_listener_test.go -package=main net Listener
//go:generate mockgen -destination=./mock_net_addr_test.go -package=main net Addr
//go:generate mockgen -destination=./mock_conn_test.go --package=main net Conn
//go:generate mockgen -destination=./mock_dir_entry_test.go --package=main github.com/Rmarken5/file-broadcaster/file-listener IFileListener
//go:generate mockgen -destination=./mock_subject_test.go -package=main github.com/Rmarken5/file-broadcaster/observer Subject

type server struct {
	FileListener file_listener.IFileListener
	FileSubject  observer.Subject
}

var directory = flag.String("directory", "", "Used to set directory of file to listen to.")
var serverPort = flag.String("port", "8000", "Used to set port of server to listen to incoming connections.")
var help = flag.Bool("help", false, "Print this menu.")
var serverAddress string

func init() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	directory := strings.TrimSpace(*directory)
	if directory == "" {
		fmt.Fprint(os.Stderr, "directory flag is required\n")
		os.Exit(99)
	}

	serverPort := strings.TrimSpace(*serverPort)
	if serverPort == "" {
		fmt.Fprint(os.Stderr, "port flag is required\n")
		os.Exit(99)
	}
	serverAddress = fmt.Sprintf("%s:%s", "0.0.0.0", serverPort)
}

func main() {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		panic(err)
	}
	dirEntries, err := os.ReadDir(*directory)
	if err != nil {
		panic(err)
	}
	s := server{
		FileListener: &file_listener.FileListener{
			Watcher: watcher,
		},
		FileSubject: &observer.FileBroadcastSubject{
			Files:     []common.FileInfo{},
			Observers: make(map[string]observer.Observer, 0),
		},
	}
	done := make(chan bool)
	a, err := net.ResolveTCPAddr("tcp", serverAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := net.ListenTCP("tcp", a)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	s.addFilesToSubject(*directory, dirEntries)
	go s.acceptClients(l)
	go s.listenForFiles(*directory)

	for {
		select {
		case <-done:
			os.Exit(0)
		}
	}

}

func (s *server) acceptClients(listener net.Listener) {
	rand.Seed(time.Now().Unix())

	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		if c != nil {
			go s.handleConnection(c)
		}
	}
}

func (s *server) handleConnection(c net.Conn) {
	buffer := make([]byte, 1024)
	var clientConnType observer.ConnHandlerMessages
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	r, err := c.Read(buffer)
	fmt.Printf("Read length: %d\n", r)
	if err != nil {
		fmt.Printf("error reading incoming connection for type: %v\n", err)
		return
	}

	if r <= 0 {
		fmt.Println("No bytes read from client on handshake")
	}
	str := strings.TrimSpace(string(buffer[0:r]))
	connType, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("error converting bytes to string to int: %v\n", "err")
		return
	}

	if !observer.IsCHM(connType) {
		fmt.Printf("Not a request that can be fulfilled: %d\n", connType)
		return
	}
	clientConnType = observer.ConnHandlerMessages(connType)
	fmt.Println("Conn handler message: " + clientConnType.String())

	if clientConnType == observer.FILE_LISTENER_CONN_TYPE {
		obs := &observer.ConnectionData{
			Address: c.RemoteAddr().String(),
			Conn:    c,
		}
		fmt.Println("Addr: ", obs.GetIdentifier())

		s.FileSubject.Subscribe(obs)
	}

	if clientConnType == observer.FILE_REQUEST_CONN_TYPE {
		if err := serveFile(c); err != nil {
			fmt.Printf("error serving file: %v", err)
		}
	}
}

func serveFile(c net.Conn) error {
	buffer := make([]byte, 1024)
	_, err := fmt.Fprintln(c, observer.SERVER_READY_TO_RECEIVE_FILE_REQUEST)
	if err != nil {
		return fmt.Errorf("not able to communicate with client: %v\n", err)
	}

	r, err := c.Read(buffer)
	if err != nil {
		return fmt.Errorf("error reading incoming connection for type: %v\n", err)
	}

	if r <= 0 {
		return fmt.Errorf("no bytes read from client")
	}
	str := strings.TrimSpace(string(buffer[0:r]))

	fmt.Println("Client wants: " + str)

	f, err := os.Open(*directory + "/" + str)
	defer f.Close()

	if err != nil {
		return fmt.Errorf("not able to open file: %v\n", err)
	}

	reader := bufio.NewReader(f)
	writer := bufio.NewWriter(c)
	defer writer.Flush()
	_, err = io.Copy(writer, reader)

	if err != nil {
		return fmt.Errorf("Unable to copy file to connection: %v\n", err)
	}
	return nil
}

func (s *server) listenForFiles(directory string) error {

	fileListener := s.FileListener.ListenForFiles(directory)
	fmt.Println("listening for files.")

	done := make(chan bool)
	go s.evaluateEvent(fileListener)

	<-done
	return nil
}

func (s *server) addFilesToSubject(dir string, dirEntries []os.DirEntry) {
	fullPath := make([]string, 0)
	files := s.FileListener.ReadDirectory(dirEntries)
	for _, file := range files {
		fullPath = append(fullPath, dir+"/"+file)
	}
	fileInfos := file_listener.BuildFileInfosFromPaths(fullPath)
	s.FileSubject.SetFiles(append(s.FileSubject.GetFiles(), fileInfos...))
}

func (s *server) evaluateEvent(listenerChannel <-chan fsnotify.Event) {
	waitFor := 100 * time.Millisecond
	timers := make(map[string]*time.Timer)
	var mu sync.Mutex

	// watch for events
	for event := range listenerChannel {
		fmt.Printf("EVENT! %+v\n", event)
		fileName := file_listener.GetFileNameFromPath(event.Name)

		if strings.HasSuffix(fileName, "~") ||
			!(event.Has(fsnotify.Create) ||
				event.Has(fsnotify.Write) ||
				event.Has(fsnotify.Remove)) {
			continue
		}
		fileInfo, err := file_listener.BuildFileInfoFromPath(event.Name)
		if err != nil {
			fmt.Printf("cannot get fileInfo for %s\n", event.Name)
			continue
		}

		// Get timer.
		mu.Lock()
		t, ok := timers[fileInfo.FileName]
		mu.Unlock()

		// No timer yet, so create one.
		if !ok {
			if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
				t = time.AfterFunc(math.MaxInt64, func() { s.FileSubject.AddFile(fileInfo) })
			} else if event.Has(fsnotify.Remove) {
				t = time.AfterFunc(math.MaxInt64, func() { s.FileSubject.RemoveFile(fileInfo) })
			}
			t.Stop()

			mu.Lock()
			timers[event.Name] = t
			mu.Unlock()
		}
		t.Reset(waitFor)
	}
}
