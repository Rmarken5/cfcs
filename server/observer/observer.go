package observer

import (
	"fmt"
	"net"
	"strings"
)

//go:generate mockgen -destination=./mock_conn_test.go --package=observer net Conn

type ConnectionData struct {
	Address string
	Conn    net.Conn
}

func (c *ConnectionData) LoadAllFiles(files []string) error {
	fileString := strings.Join(files, ",") + "\n"

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
