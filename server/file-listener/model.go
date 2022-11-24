package file_listener

import (
	"github.com/fsnotify/fsnotify"
	"os"
)

//go:generate mockgen -destination=./mock_file_listener_test.go -package=file_listener . IFileListener


type IFileListener interface {
	ListenForFiles(directory string) chan fsnotify.Event
	ReadDirectory(dirEntries []os.DirEntry) []string
}

type FileListener struct {
	Watcher *fsnotify.Watcher
}
