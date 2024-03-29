package main

import (
	"flag"
	"fmt"
	"github.com/rmarken5/cfcs/client/client"
	file_manager "github.com/rmarken5/cfcs/client/file-manager"
	"os"
	"strings"
)

var serverHost = flag.String("host", "localhost", "Used to set host of server to connect to.")
var serverPort = flag.String("port", "8000", "Used to set port of server to connect to.")
var fileDestination = flag.String("directory", "/tmp", "Required. Used to set destination of downloaded file")
var dbLocation = flag.String("db-location", "./db", "Used to set the location of the database that tracks what files have been downloaded")
var help = flag.Bool("help", false, "Print this menu.")
var serverAddress string
var db string
var fileDest string

func init() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	fileDest = strings.TrimSpace(*fileDestination)
	if fileDest == "" {
		fmt.Fprint(os.Stderr, "directory flag is required\n")
		os.Exit(99)
	}

	serverHost := strings.TrimSpace(*serverHost)
	if serverHost == "" {
		fmt.Fprint(os.Stderr, "host is required\n")
		flag.PrintDefaults()
		os.Exit(99)
	}

	serverPort := strings.TrimSpace(*serverPort)
	if serverPort == "" {
		fmt.Fprint(os.Stderr, "port is required\n")
		flag.PrintDefaults()
		os.Exit(99)
	}

	db = strings.TrimSpace(*dbLocation)
	if db == "" {
		fmt.Fprint(os.Stderr, "db-location is required\n")
		flag.PrintDefaults()
		os.Exit(99)
	}

	if _, err := os.Stat(db); os.IsNotExist(err) {
		if err := os.MkdirAll(db, os.FileMode(0755)); err != nil {
			fmt.Errorf("Failed to create dir: %w\n", err)
		}
	}

	serverAddress = fmt.Sprintf("%s:%s", serverHost, serverPort)
}

func main() {

	fileManagerImpl := file_manager.NewFileManagerImpl(db)
	clientImpl, err := client.NewClientImpl(fileManagerImpl, fileDest)
	if err != nil {
		panic(err)
	}

	conn, err := clientImpl.ConnectToServer(serverAddress)
	if err != nil {
		panic(err)
	}

	go clientImpl.ManageServerResponses(conn)
	go clientImpl.RequestFiles(serverAddress)

	err = clientImpl.RequestAllFileNames(conn)
	if err != nil {
		panic(err)
	}

	forever := make(chan int)
	<-forever
}
