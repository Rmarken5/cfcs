package observer

import (
	"fmt"
	"github.com/rmarken5/cfcs/common"
	"sync"
)

type FileBroadcastSubject struct {
	Files     []common.FileInfo
	Observers map[string]Observer
	fileMutex sync.RWMutex
}

func (f *FileBroadcastSubject) AddFile(fileInfo common.FileInfo) {
	var isExists bool
	for _, file := range f.Files {
		if file.Eq(&fileInfo) {
			isExists = true
		}
	}
	if !isExists {
		f.Files = append(f.Files, fileInfo)
	}
	f.NotifyAllWithFile(fileInfo)
}
func (f *FileBroadcastSubject) RemoveFile(fileInfo common.FileInfo) {
	newFileArr := f.Files
	for i, file := range f.Files {
		if file.Eq(&fileInfo) {
			newFileArr = append(f.Files[:i], f.Files[i+1:]...)
		}
	}
	f.Files = newFileArr
}

func (f *FileBroadcastSubject) Subscribe(observer Observer) {
	fmt.Println("Adding new observer to subject: ", observer.GetIdentifier())
	f.Observers[observer.GetIdentifier()] = observer
	observer.LoadAllFiles(f.Files)
}

func (f *FileBroadcastSubject) Unsubscribe(key string) {
	delete(f.Observers, key)
	fmt.Printf("%s has closed their connection.\n", key)
}

func (f *FileBroadcastSubject) NotifyAllWithFiles(files []common.FileInfo) {
	for _, obs := range f.Observers {
		if err := obs.LoadAllFiles(files); err != nil {
			fmt.Printf("Connection closed for: %s\n", obs.GetIdentifier())
			f.Unsubscribe(obs.GetIdentifier())
		}
	}
}

func (f *FileBroadcastSubject) NotifyAllWithFile(fileInfo common.FileInfo) {
	for _, obs := range f.Observers {
		if err := obs.AddFile(fileInfo); err != nil {
			fmt.Printf("Connection closed for: %s\n", obs.GetIdentifier())
			f.Unsubscribe(obs.GetIdentifier())
		}
	}
}

func (f *FileBroadcastSubject) SetFiles(files []common.FileInfo) {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()
	f.Files = files
}

func (f *FileBroadcastSubject) GetFiles() []common.FileInfo {
	f.fileMutex.Lock()
	defer f.fileMutex.Unlock()
	return f.Files
}
