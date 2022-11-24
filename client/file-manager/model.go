package file_manager

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//go:generate mockgen -destination file_manager_test/mock_file_manager_test.go -package file_manager_test net Conn

type FileManager interface {
	WriteSrcToDest(srcFile *sftp.File, fullyQualifiedPath string) (string, error)
	GetSourceFile(sshConfig *ssh.ClientConfig, fullyQualifiedPath string) (*sftp.File, error)
	WriteFileHashToDB(fileName string, file *sftp.File) error
	RemoveFileFromQueue(fileName string)
	CloseConns()
	ShouldWriteToDB(fileName string, srcFile *sftp.File) bool
}
