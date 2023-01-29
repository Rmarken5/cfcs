package observer

import (
	"encoding/json"
	"fmt"
	"github.com/rmarken5/cfcs/common"
	"net"
)

//go:generate mockgen -destination=./mock_conn.go --package=observer net Conn

type ConnectionObserver struct {
	Address string
	Conn    net.Conn
}

func (c *ConnectionObserver) LoadAllFiles(files []common.FileInfo) error {

	if _, err := c.Conn.Write([]byte(fmt.Sprintf("%d\n", common.SERVER_SENDING_FILE_LIST))); err != nil {
		fmt.Printf("Unable to write %s to %s\n", common.SERVER_SENDING_FILE_LIST, c.Address)
		return fmt.Errorf("error %v\n: ", err)
	}

	fmt.Printf("writing files: %+v\n", files)
	for _, file := range files {
		if err := c.AddFile(file); err != nil {
			return err
		}
	}

	fmt.Println("Files written.")
	return nil
}

func (c *ConnectionObserver) AddFile(file common.FileInfo) error {
	fmt.Printf("writing file: %v\n", file)

	if err := json.NewEncoder(c.Conn).Encode(file); err != nil {
		return fmt.Errorf("Unable to write %v to %s: %v\n", file, c.Address, err)
	}
	fmt.Println("File written.")
	return nil
}

func (c *ConnectionObserver) GetIdentifier() string {
	return c.Address
}

