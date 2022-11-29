package client

import (
	"net"
)

//go:generate mockgen -destination mock_client_test.go -package client_test net Conn

//go:generate mockgen -destination mock_client_test.go -package client_test github.com/Rmarken5/file-client/file-manager FileManager

type Client interface {
	ConnectToServer(address string) (*net.TCPConn, error)
	RequestAllFileNames(conn *net.TCPConn) error
	ManageServerResponses(conn *net.TCPConn)
	ListenForFiles(conn *net.TCPConn)
	RequestFiles(serverAddress string)
	WriteComplete(conn *net.TCPConn)
	FileChannel() *chan string
	DownloadFiles()
}
