package client

import (
	"net"
)

//go:generate mockgen -destination mock_client.go -package client . Client

//go:generate mockgen -destination mock_net_conn.go -package client net Conn

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
