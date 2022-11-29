package file_manager

import (
	"os"
)

//go:generate mockgen -destination conn_manager_test/mock_conn_manager_test.go -package file_manager_test net Conn

type ConnectionManager interface {
	HandleServerResponse(response string) error
	GetAllFileNamesFromServer() ([]string, error)
	RequestFileFromServer(fileName string) (*os.File, error)
	WriteFileHashToDB(fileName string, file *os.File) error
	RemoveFileFromQueue(fileName string)
	CloseConns()
	ShouldWriteToDB(fileName string, srcFile *os.File) bool
}
