package observer

import (
	"fmt"
	"net"
	"strings"
)

//go:generate mockgen -destination=./mock_conn_test.go --package=observer net Conn

type ConnHandlerMessages int

const (
	FILE_LISTENER_CONN_TYPE ConnHandlerMessages = iota
	FILE_REQUEST_CONN_TYPE
	SERVER_READY_TO_RECIEVE_FILE_REQUEST
	SERVER_SENDING_FILE_LIST
)

type ConnectionData struct {
	Address           string
	Conn              net.Conn
	ClientRequestType ConnHandlerMessages
}

func (c *ConnectionData) LoadAllFiles(files []string) error {
	fileString := strings.Join(files, ",") + "\n"

	if _, err := c.Conn.Write([]byte(fmt.Sprintf("%d\n", SERVER_SENDING_FILE_LIST))); err != nil {
		fmt.Printf("Unable to write %s to %s", fileString, c.Address)
		return fmt.Errorf("error %v: ", err)
	}

	fmt.Println("writing files: ", fileString)

	if _, err := c.Conn.Write([]byte(fileString)); err != nil {
		fmt.Printf("Unable to write %s to %s", fileString, c.Address)
		return fmt.Errorf("error %v: ", err)
	}
	fmt.Println("File String written.")
	return nil
}

func (c *ConnectionData) AddFile(file string) error {

	fmt.Println("writing file: ", file)

	if _, err := c.Conn.Write([]byte(file + "\n")); err != nil {
		fmt.Printf("Unable to write %s to %s", file, c.Address)
		return fmt.Errorf("error %v: ", err)
	}
	fmt.Println("File String written.")
	return nil
}

func (c *ConnectionData) GetIdentifier() string {
	return c.Address
}

func IsCCT(n int) bool {
	conv := ConnHandlerMessages(n)
	return conv == FILE_LISTENER_CONN_TYPE || conv == FILE_REQUEST_CONN_TYPE
}

func (chm ConnHandlerMessages) String() string {
	switch chm {
	case FILE_LISTENER_CONN_TYPE:
		return "FILE_LISTENER_CONNECTION"
	case FILE_REQUEST_CONN_TYPE:
		return "FILE_REQUEST_CONNECTION"
	case SERVER_READY_TO_RECIEVE_FILE_REQUEST:
		return "SERVER_READY_TO_RECEIVE_FILE_REQUEST"
	}
	return ""
}
