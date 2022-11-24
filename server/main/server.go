package main

import (
	"flag"
	"fmt"
	"github.com/Rmarken5/file-broadcaster/observer"
	"github.com/fsnotify/fsnotify"
	file_listener "github.com/rmarken5/cfcs/server/file-listener"
	"math/rand"
	"net"
	"os"
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

	l, err := net.Listen("tcp4", ":8000")
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
		go s.handleConnection(c)
	}
}

func (s *server) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	obs := &observer.ConnectionData{
		Address: c.RemoteAddr().String(),
		Conn:    c,
	}
	fmt.Println("Addr: ", obs.GetIdentifier())
	s.FileSubject.Subscribe(obs)
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
