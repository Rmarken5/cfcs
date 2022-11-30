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
	//buffer := make([]byte, 1024)
	// TODO: Create tcp connection from Config file or cl options
	flag.Parse()

	fileManagerImpl := file_manager.NewConnectionManagerImpl("/home/ryan/file-client/")
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

	/*write, err := fmt.Fprintf(tcp, "%d\n", observer.FILE_REQUEST_CONN_TYPE)

	if err != nil {
		fmt.Printf("error writing to tcp: %v\n", err)
	}
	fmt.Printf("written: %d\n", write)

	n, err := tcp.Read(buffer)
	if err != nil {
		fmt.Printf("error reading from connection: %v\n", err)
	}

	str := strings.TrimSpace(string(buffer[0:n]))
	fmt.Println(str)

	write, err = fmt.Fprintf(tcp, "%s\n", "test.txt")

	f, err := os.Create("/home/ryan/programming/go-programs/cfcs/client/tmp/test.txt")
	defer f.Close()
	if err != nil {
		fmt.Printf("not able to open file: %v\n", err)
		return
	}
	reader := bufio.NewReader(tcp)
	writer := bufio.NewWriter(f)
	defer writer.Flush()

	_, err = io.Copy(writer, reader)
	if err != nil {
		fmt.Printf("Unable to copy file to connection: %v\n", err)
		return
	}*/

	/*
		sshClient, err := ssh.Dial("tcp", ":22", &clientConfig)
		if err != nil {
			panic("Failed to dial: " + err.Error())
		}
		fmt.Println("Successfully connected to ssh server.")

		// open an SFTP session over an existing ssh connection.
		sftp, err := sftp.NewClient(sshClient)
		if err != nil {
			log.Fatal(err)
		}

		fileManagerImpl := file_manager.NewConnectionManagerImpl("/home/ryan/file-client/", sftp)
		clientImpl := client.NewClientImpl(fileManagerImpl)

		conn, err := clientImpl.ConnectToServer(*serverAddress)
		if err != nil {
			log.Fatalf("Cannot open a connection to %s: %v\n", *serverAddress, err)
		}

		go clientImpl.ListenForFiles(conn)
		go gracefulShutdown(fileManagerImpl)
		go clientImpl.DownloadFiles(clientConfig)
	*/
	//fileManagerImpl.PrintKeysInDatabase()

	/*	go func() {
		for {
			select {
			case file := <-*clientImpl.FileChannel():
				go func() {
					err := fileManagerImpl.DownloadFile(&clientConfig, "/home/ryanm/programming/go/file-broadcaster/files/"+file)
					if err != nil {
						fmt.Println(err)
					}
				}()
			}
		}
	}()*/
	forever := make(chan int)
	<-forever
}

func gracefulShutdown(fileManager file_manager.FileManager) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		fmt.Println("Sutting down gracefully.")
		fileManager.CloseConns()
		os.Exit(0)
	}()
}
