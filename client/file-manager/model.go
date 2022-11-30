package file_manager

import (
	"github.com/rmarken5/cfcs/common"
)

//go:generate mockgen -destination file_manager_test/mock_file_manager_test.go -package file_manager_test net Conn

type FileManager interface {
	WriteFileHashToDB(fileInfo common.FileInfo) error
	RemoveFileFromQueue(fileName string)
	CloseConns()
	ShouldWriteToDB(fileInfo common.FileInfo) bool
}
