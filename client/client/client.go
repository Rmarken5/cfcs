package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	file_manager "github.com/rmarken5/cfcs/client/file-manager"
	"github.com/rmarken5/cfcs/common"
	"github.com/rmarken5/cfcs/server/observer"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type ClientImpl struct {
	fileManager *file_manager.FileManagerImpl
	fileChannel chan common.FileInfo
	directory string
}

func NewClientImpl(fileManager *file_manager.FileManagerImpl, directory string) *ClientImpl {
	fileChannel := make(chan common.FileInfo)
	return &ClientImpl{
		fileManager: fileManager,
		fileChannel: fileChannel,
		directory: directory,
	}
}

func (c *ClientImpl) ConnectToServer(address string) (*net.TCPConn, error) {

	rAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		fmt.Printf("error connecting to tcp: %v\n", err)
	}
	tcp, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		fmt.Printf("error connecting to tcp: %v\n", err)
	}
	return tcp, err
}

func (c *ClientImpl) RequestAllFileNames(conn *net.TCPConn) error {
	w, err := fmt.Fprintf(conn, "%d\n", observer.FILE_LISTENER_CONN_TYPE)
	if err != nil {
		return fmt.Errorf("error requesting all files: %v", err)
	}
	fmt.Printf("Requesing file name bytes writen: %d\n", w)
	return nil
}

func (c *ClientImpl) ManageServerResponses(conn *net.TCPConn) {
	fmt.Println("Listening on reads from server")
	for {
		buf := make([]byte, 1024)
		var str string

		n, err := conn.Read(buf)

		if err != nil {
			if err != io.ErrUnexpectedEOF {
				fmt.Fprintln(os.Stderr, err)
			}
		}

		fmt.Println("read n bytes...", n)
		// process buf
		str += string(buf[0:n])

		str = strings.TrimSpace(str)
		fmt.Println("String from server: " + str)
		time.Sleep(time.Second)
		switch str {
		case fmt.Sprintf("%d", observer.SERVER_SENDING_FILE_LIST):
			c.ListenForFiles(conn)
		}
	}
}

func (c *ClientImpl) ListenForFiles(conn *net.TCPConn) {
	fmt.Println("Starting to listen to files.")
	reader := bufio.NewReader(conn)
	for {
		fmt.Println("Starting loop")
		readString, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		fmt.Println(readString)
		if len(readString) > 0 {
			var info common.FileInfo
			err = json.Unmarshal([]byte(readString), &info)
			if err != nil {
				fmt.Printf("cannot unmarshal bytes: %v\n", err)
				continue
			}
			fmt.Printf("File: %+v\n", info)
			c.fileChannel <- info
			fmt.Printf("received %s\n", info)
		}
	}
}

func (c *ClientImpl) RequestFiles(serverAddress string) {
	for file := range c.fileChannel {
		if c.fileManager.ShouldWriteToDB(file) {
			go func(info common.FileInfo) {
				fmt.Printf("got %v\n", info)
				buffer := make([]byte, 1024)
				tcp, err := c.ConnectToServer(serverAddress)
				defer func(tcp *net.TCPConn) {
					err := tcp.Close()
					if err != nil {
						fmt.Printf("error closing conn: %v\n", err)
					}
				}(tcp)
				if err != nil {
					fmt.Printf("Cannot conenct to server: %v\n", err)
					return
				}
				write, err := fmt.Fprintf(tcp, "%d\n", observer.FILE_REQUEST_CONN_TYPE)

				if err != nil {
					fmt.Printf("error writing to tcp: %v\n", err)
				}
				fmt.Printf("written: %d\n", write)

				n, err := tcp.Read(buffer)
				if err != nil {
					fmt.Printf("error reading from connection: %v\n", err)
				}

				str := strings.TrimSpace(string(buffer[0:n]))
				fmt.Println("Should be \"SERVER_READY\": " + str)
				f, err := os.Create(c.directory + "/" + info.FileName)
				defer f.Close()

				write, err = fmt.Fprintf(tcp, "%s\n", info.FileName)


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
				}
				err = c.fileManager.WriteFileHashToDB(file)
				if err != nil {
					fmt.Printf("Error Downloading file: %s. Error: %v", file, err)
				}
			}(file)
		}
	}
}
func (c *ClientImpl) FileChannel() *chan common.FileInfo {
	return &c.fileChannel
}
