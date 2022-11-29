package client

import (
	"bufio"
	"bytes"
	"fmt"
	file_manager "github.com/rmarken5/cfcs/client/file-manager"
	"github.com/rmarken5/cfcs/server/observer"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type ClientImpl struct {
	fileManager *file_manager.ConnectionManagerImpl
	fileChannel chan string
}

func NewClientImpl(fileManager *file_manager.ConnectionManagerImpl) *ClientImpl {
	fileChannel := make(chan string)
	return &ClientImpl{
		fileManager: fileManager,
		fileChannel: fileChannel,
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

		readString, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		fmt.Println(readString)
		files := strings.Split(readString, ",")
		fmt.Printf("Files: %+v\n", files)
		fmt.Println(len(files))

		for _, file := range files {
			file = strings.TrimSpace(file)
			file = string(bytes.Trim([]byte(file), "\x00"))
			c.fileChannel <- file
			fmt.Printf("received %s\n", file)
		}
	}
}

func (c *ClientImpl) RequestFiles(serverAddress string) {
	for file := range c.fileChannel {
		go func(fileName string) {
			fmt.Printf("got %s", fileName)
			buffer := make([]byte, 1024)
			tcp, err := c.ConnectToServer(serverAddress)
			defer tcp.Close()
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
			fmt.Println(str)

			write, err = fmt.Fprintf(tcp, "%s\n", fileName)

			f, err := os.Create("/home/ryan/programming/go-programs/cfcs/client/tmp/" + fileName)
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
			}
			err = c.fileManager.WriteFileHashToDB(fileName, f)
			if err != nil {
				fmt.Printf("Error Downloading file: %s. Error: %v", file, err)
			}
		}(file)
	}
}
func (c *ClientImpl) WriteComplete(conn *net.Conn) {
	panic("implement me")
}

func (c *ClientImpl) FileChannel() *chan string {
	return &c.fileChannel
}

// DownloadFiles listens to the file channel for incoming files and writes them if they don't exist in DB.
/*func (c *ClientImpl) DownloadFiles(config ssh.ClientConfig) {
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
}*/
