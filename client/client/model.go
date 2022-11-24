package client

import (
	"net"
)

//go:generate mockgen -destination mock_client_test.go -package client_test net Conn

//go:generate mockgen -destination mock_client_test.go -package client_test github.com/Rmarken5/file-client/file-manager FileManager

type Client interface {
	ConnectToServer(address string) (net.Conn, error)
	ListenForFiles(conn *net.Conn, fileChannel <-chan string)
	WriteComplete(conn *net.Conn)
	FileChannel() *chan string
	DownloadFiles()
}
