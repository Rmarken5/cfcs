package observer

import (
	"encoding/json"
	"fmt"
	"github.com/rmarken5/cfcs/common"
	"net"
)

//go:generate mockgen -destination=./mock_conn_test.go --package=observer net Conn

type ConnHandlerMessages int

const (
	FILE_LISTENER_CONN_TYPE ConnHandlerMessages = iota
	FILE_REQUEST_CONN_TYPE
	SERVER_READY_TO_RECEIVE_FILE_REQUEST
	SERVER_SENDING_FILE_LIST
)

type ConnectionData struct {
	Address string
	Conn    net.Conn
}

func (c *ConnectionData) LoadAllFiles(files []common.FileInfo) error {

	if _, err := c.Conn.Write([]byte(fmt.Sprintf("%d\n", SERVER_SENDING_FILE_LIST))); err != nil {
		fmt.Printf("Unable to write %s to %s", SERVER_SENDING_FILE_LIST, c.Address)
		return fmt.Errorf("error %v: ", err)
	}

	fmt.Printf("writing files: %v+", files)
	for _, file := range files {
		if err := json.NewEncoder(c.Conn).Encode(file); err != nil {
			return fmt.Errorf("Unable to write %v to %s: %v\n", files, c.Address, err)
		}
	}

	fmt.Println("Files written.")
	return nil
}

func (c *ConnectionData) AddFile(file common.FileInfo) error {

	fmt.Printf("writing file: %v\n", file)

	if err := json.NewEncoder(c.Conn).Encode(file); err != nil {
		return fmt.Errorf("Unable to write %v to %s: %v\n", file, c.Address, err)
	}
	fmt.Println("File written.")
	return nil
}

func (c *ConnectionData) GetIdentifier() string {
	return c.Address
}

func IsCHM(n int) bool {
	conv := ConnHandlerMessages(n)
	return conv == FILE_LISTENER_CONN_TYPE || conv == FILE_REQUEST_CONN_TYPE
}

func (chm ConnHandlerMessages) String() string {
	switch chm {
	case FILE_LISTENER_CONN_TYPE:
		return "FILE_LISTENER_CONNECTION"
	case FILE_REQUEST_CONN_TYPE:
		return "FILE_REQUEST_CONNECTION"
	case SERVER_READY_TO_RECEIVE_FILE_REQUEST:
		return "SERVER_READY_TO_RECEIVE_FILE_REQUEST"
	}
	return ""
}
