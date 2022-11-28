package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	file_listener "github.com/rmarken5/cfcs/server/file-listener"
	"github.com/rmarken5/cfcs/server/observer"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
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

var directory = flag.String("directory", "../dummy", "Directory to listen to files on.")

func main() {

	flag.Parse()

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
			Files:     []string{},
			Observers: make(map[string]observer.Observer, 0),
		},
	}
	done := make(chan bool)
	SERVER := "localhost" + ":" + "8000"
	a, err := net.ResolveTCPAddr("tcp", SERVER)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := net.ListenTCP("tcp4", a)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	s.addFilesToSubject(dirEntries)
	go s.acceptClients(l)
	go s.listenForFiles(*directory)

	for {
		select {
		case <-done:
			os.Exit(1)
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
	var str string
	for {
		r, err := c.Read(buffer)
		fmt.Println(r)
		if err != nil {
			fmt.Printf("error reading incoming connection for type: %v\n", err)
			return
		}

		if r > 0 {
			str = strings.TrimSpace(string(buffer[0:r]))
			fmt.Println(str)
			break
		}
	}

	connType, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("error converting bytes to string to int: %v\n", "err")
		return
	}

	if !observer.IsCCT(connType) {
		fmt.Printf("Not a request that can be fulfilled: %d\n", connType)
		return
	}
	clientConnType = observer.ConnHandlerMessages(connType)
	fmt.Println(clientConnType.String())

	obs := &observer.ConnectionData{
		Address:           c.RemoteAddr().String(),
		Conn:              c,
		ClientRequestType: clientConnType,
	}
	fmt.Println("Addr: ", obs.GetIdentifier())
	if clientConnType == observer.FILE_LISTENER_CONN_TYPE {
		s.FileSubject.Subscribe(obs)
	}

	if clientConnType == observer.FILE_REQUEST_CONN_TYPE {
		_, err = fmt.Fprintln(c, observer.SERVER_READY_TO_RECIEVE_FILE_REQUEST)
		if err != nil {
			fmt.Printf("not able to communicate with client: %v\n", err)
			return
		}
		for {
			r, err := c.Read(buffer)
			fmt.Printf("Read length: %d/n", r)
			if err != nil {
				fmt.Printf("error reading incoming connection for type: %v\n", err)
				return
			}

			if r > 0 {
				str = strings.TrimSpace(string(buffer[0:r]))
				break
			}
		}

		f, err := os.Open(*directory + "/" + str)
		defer f.Close()
		if err != nil {
			fmt.Printf("not able to open file: %v\n", err)
			return
		}
		reader := bufio.NewReader(f)
		writer := bufio.NewWriter(c)
		defer writer.Flush()

		_, err = io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("Unable to copy file to connection: %v\n", err)
			return
		}
	}
}

func (s *server) listenForFiles(directory string) error {

	fileListener := s.FileListener.ListenForFiles(directory)
	fmt.Println("listening for files.")

	done := make(chan bool)

	go func() {
		for {
			s.evaluateEvent(fileListener)
		}
	}()
	<-done
	return nil
}

func (s *server) addFilesToSubject(dirEntries []os.DirEntry) {
	files := s.FileListener.ReadDirectory(dirEntries)

	s.FileSubject.SetFiles(append(s.FileSubject.GetFiles(), files...))
}

func (s *server) evaluateEvent(listenerChannel <-chan fsnotify.Event) {
	select {
	// watch for events
	case event := <-listenerChannel:
		fmt.Printf("EVENT! %+v\n", event)
		fileParts := strings.Split(event.Name, "/")
		fileName := fileParts[len(fileParts)-1]
		if event.Op.String() == "CREATE" {
			s.FileSubject.AddFile(fileName)
		} else if event.Op.String() == "REMOVE" {
			s.FileSubject.RemoveFile(fileName)
		}
	}
}
