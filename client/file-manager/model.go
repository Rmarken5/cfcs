package file_manager

import (
	"os"
)

//go:generate mockgen -destination file_manager_test/mock_file_manager_test.go -package file_manager_test net Conn

type FileManager interface {
	WriteFileHashToDB(fileName string, file *os.File) error
	RemoveFileFromQueue(fileName string)
	CloseConns()
	ShouldWriteToDB(fileName string, srcFile *os.File) bool
}
