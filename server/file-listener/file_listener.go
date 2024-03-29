package file_listener

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rmarken5/cfcs/common"
	"io"
	"os"
	"strings"
)

//go:generate mockgen -destination=./mock_dir_entry.go --package=file_listener io/fs DirEntry
//go:generate mockgen -destination=./mock_file_listener.go -package=file_listener . FileListener

type FileListener interface {
	ListenForFiles(directory string) (chan fsnotify.Event, error)
	ReadDirectory(dirEntries []os.DirEntry) []string
	CreateDirectory(directory string) error
}

type FileListenerImpl struct {
	Watcher *fsnotify.Watcher
}

func (f *FileListenerImpl) ListenForFiles(directory string) (chan fsnotify.Event, error) {
	err := f.Watcher.Add(directory)
	if err != nil {
		return nil, fmt.Errorf("error getting event channel %w\n", err)
	}
	return f.Watcher.Events, nil
}

// ReadDirectory gets file name from the os.DirEntry - excluding entries that are directories.
func (f *FileListenerImpl) ReadDirectory(dirEntries []os.DirEntry) []string {
	var files []string
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}

func BuildFileInfoFromPath(filePath string) (common.FileInfo, error) {
	fmt.Printf("BuildFileInfoFromPath: %s\n", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		return common.FileInfo{}, fmt.Errorf("error reading file: %v", err)
	}

	hash, err := generateHash(f)
	if err != nil {
		return common.FileInfo{}, fmt.Errorf("error generating hash: %v", err)
	}
	return common.FileInfo{
		Hash:     hash,
		FileName: GetFileNameFromPath(f.Name()),
	}, nil

}

func (f *FileListenerImpl) CreateDirectory(directory string) error  {

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0766)
		return err
	}
	return nil
}

func generateHash(srcFile *os.File) (string, error) {

	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string
	//Open a new hash interface to write to
	hash := md5.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, srcFile); err != nil {
		fmt.Printf("error in copy: %v\n", err)
		return returnMD5String, err
	}
	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func GetFileNameFromPath(path string) string {
	fileParts := strings.Split(path, "/")
	fileName := fileParts[len(fileParts)-1]
	return fileName
}
