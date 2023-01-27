package file_manager

import (
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/rmarken5/cfcs/common"
	"log"
)

type FileManagerImpl struct {
	fileQueue map[string]bool
	db        *bitcask.Bitcask
}

func NewFileManagerImpl(dbLocation string) *FileManagerImpl {
	open, err := bitcask.Open(dbLocation, bitcask.WithSync(true), bitcask.WithAutoRecovery(true))
	if err != nil {
		log.Fatalf("Error opening db: %v\n", err)
	}

	return &FileManagerImpl{
		fileQueue: make(map[string]bool),
		db:        open,
	}
}

func (f *FileManagerImpl) WriteFileHashToDB(info common.FileInfo) error {
	err := f.db.Put([]byte(info.FileName), []byte(info.Hash))
	fmt.Println(f.db.Get([]byte(info.FileName)))
	return err
}

func (f *FileManagerImpl) PrintKeysInDatabase() {
	k := <-f.db.Keys()
	fmt.Println(string(k))
}

func (f *FileManagerImpl) RemoveFileFromQueue(fileName string) {
	delete(f.fileQueue, fileName)
}

func (f *FileManagerImpl) CloseConns() {
	fmt.Println("Database close called.")
	f.db.Close()
}

func (f *FileManagerImpl) ShouldWriteToDB(info common.FileInfo) bool {

	fmt.Printf("filename lookup: %v\n", info)
	if !f.db.Has([]byte(info.FileName)) {
		return true
	}
	fmt.Println("db has file")

	fmt.Println("Hash: " + info.Hash)
	hashToCompare, err := f.db.Get([]byte(info.FileName))
	if err != nil {
		fmt.Printf("error getting filename from db: %s\n %v\n", info.FileName, err)
	}
	fmt.Println("File hash: " + string(hashToCompare))
	return info.Hash != string(hashToCompare)
}

