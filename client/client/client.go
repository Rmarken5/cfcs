package client

import (
	"bufio"
	"bytes"
	"fmt"
	file_manager "github.com/rmarken5/cfcs/client/file-manager"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
)

type ClientImpl struct {
	fileManager *file_manager.ConnectionManagerImpl
	fileChannel *chan string
}

func NewClientImpl(fileManager *file_manager.ConnectionManagerImpl) *ClientImpl {
	fileChannel := make(chan string)
	return &ClientImpl{
		fileManager: fileManager,
		fileChannel: &fileChannel,
	}
}

func (c *ClientImpl) ConnectToServer(address string) (net.Conn, error) {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to server: %v\n", err)
	}
	log.Printf("Connection Opened to %s\n", address)

	return conn, err

}

func (c *ClientImpl) ListenForFiles(conn net.Conn) {
	fmt.Println("Starting to listen to files.")
	reader := bufio.NewReader(conn)
	for {

		file, err := reader.ReadString('\n')

		if err != nil {
			log.Println(err)
		}
		files := strings.Split(file, ",")
		fmt.Printf("Files: %+v", files)
		fmt.Println(len(files))

		for _, file := range files {
			file = strings.TrimSpace(file)
			file = string(bytes.Trim([]byte(file), "\x00"))

			fmt.Printf("received %s\n", file)
			*c.fileChannel <- file
		}
	}
}

func (c *ClientImpl) WriteComplete(conn *net.Conn) {
	panic("implement me")
}

func (c *ClientImpl) FileChannel() *chan string {
	return c.fileChannel
}

// DownloadFiles listens to the file channel for incoming files and writes them if they don't exist in DB.
func (c *ClientImpl) DownloadFiles(config ssh.ClientConfig) {
	go func() {
		for file := range *c.fileChannel {
			// DownloadFiles TODO: Refactor server to not need fully qualified path from client
			sourceFile, err := c.fileManager.GetSourceFile(&config, "/home/ryan/programming/go-programs/file-broadcaster/dummy/"+file)
			if err != nil {
				if sourceFile != nil {
					sourceFile.Close()
				}
				fmt.Printf("Error getting src file: %s. Error: %v", file, err)
			}
			if c.fileManager.ShouldWriteToDB(file, sourceFile) {
				fmt.Printf("Processing " + file)
				_, err := c.fileManager.WriteSrcToDest(sourceFile, "/home/ryan/programming/go-programs/file-broadcaster/dummy/"+file)
				if err != nil {
					fmt.Printf("Error writing dest file: %s. Error: %v", file, err)
					sourceFile.Close()
				}
				err = c.fileManager.WriteFileHashToDB(file, sourceFile)
				if err != nil {
					fmt.Printf("Error Downloading file: %s. Error: %v", file, err)
					sourceFile.Close()
				}
				sourceFile.Close()
			}
		}
	}()
}
