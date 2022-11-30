package main

import (
	"flag"
	"fmt"
	"github.com/rmarken5/cfcs/client/client"
	file_manager "github.com/rmarken5/cfcs/client/file-manager"
	"os"
	"os/signal"
	"syscall"
)

var serverAddress = flag.String("directory", "localhost:8000", "Directory to listen to files on.")

func main() {

	// TODO: Create tcp connection from Config file or cl options
	flag.Parse()

	fileManagerImpl := file_manager.NewFileManagerImpl("/home/ryan/file-client/")
	clientImpl := client.NewClientImpl(fileManagerImpl)

	conn, err := clientImpl.ConnectToServer(*serverAddress)
	if err != nil {
		panic(err)
	}
	go clientImpl.ManageServerResponses(conn)
	go clientImpl.RequestFiles(*serverAddress)

	err = clientImpl.RequestAllFileNames(conn)
	if err != nil {
		panic(err)
	}

	forever := make(chan int)
	<-forever
}

func gracefulShutdown(fileManager file_manager.FileManager) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		fmt.Println("Shutting down gracefully.")
		fileManager.CloseConns()
		os.Exit(0)
	}()
}
