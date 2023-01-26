package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rmarken5/cfcs/common"
	file_listener "github.com/rmarken5/cfcs/server/file-listener"
	"github.com/rmarken5/cfcs/server/file-server"
	"github.com/rmarken5/cfcs/server/observer"
	"net"
	"os"
	"strings"
)

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
	done := make(chan bool)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	s := file_server.Server{
		FileListener: &file_listener.FileListenerImpl{
			Watcher: watcher,
		},
		FileSubject: &observer.FileBroadcastSubject{
			Files:     []common.FileInfo{},
			Observers: make(map[string]observer.Observer, 0),
		},
		FileDirectory: *directory,
	}

	s.FileListener.CreateDirectory(*directory)

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
	defer func(l *net.TCPListener) {
		err := l.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(l)

	err = s.AddFilesToSubject(*directory)
	if err != nil {
		fmt.Printf("error adding files to subject %v\n", err)
		return
	}
	go s.AcceptClients(l)
	go func() {
		err := s.ListenForFiles(*directory)
		if err != nil {
			fmt.Printf("error listening to files %v\n", err)
			done <- true
		}
	}()

	for {
		select {
		case <-done:
			os.Exit(0)
		}
	}
}
