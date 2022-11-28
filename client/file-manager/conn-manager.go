package file_manager

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"io"
	"log"
	"os"
)

type ConnectionManagerImpl struct {
	fileQueue map[string]bool
	db        *bitcask.Bitcask
}

func NewFileManagerImpl(dbLocation string) *ConnectionManagerImpl {
	open, err := bitcask.Open(dbLocation, bitcask.WithSync(true), bitcask.WithAutoRecovery(true))
	if err != nil {
		log.Fatalf("Error opening db: %v\n", err)
	}

	return &ConnectionManagerImpl{
		fileQueue: make(map[string]bool),
		db:        open,
	}
}

func (f *ConnectionManagerImpl) WriteFileHashToDB(fileName string, srcFile *os.File) error {
	hash, err := generateHash(srcFile)
	if err != nil {
		return err
	}
	err = f.db.Put([]byte(fileName), []byte(hash))
	fmt.Println(f.db.Get([]byte(fileName)))
	return err
}

func (f *ConnectionManagerImpl) PrintKeysInDatabase() {
	k := <-f.db.Keys()
	fmt.Println(string(k))
}

func (f *ConnectionManagerImpl) RemoveFileFromQueue(fileName string) {
	delete(f.fileQueue, fileName)
}

func (f *ConnectionManagerImpl) CloseConns() {
	fmt.Println("Database close called.")
	f.db.Close()
}

func (f *ConnectionManagerImpl) ShouldWriteToDB(fileName string, srcFile *os.File) bool {

	fmt.Println("filename lookup: " + fileName)
	if !f.db.Has([]byte(fileName)) {
		return true
	}
	fmt.Println("db has file")

	hash, err := generateHash(srcFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	fmt.Println("Hash: " + hash)
	hashToCompare, err := f.db.Get([]byte(fileName))
	fmt.Println("File hash: " + string(hashToCompare))
	return hash != string(hashToCompare)
}

func (f *ConnectionManagerImpl) HandleServerResponse(response string) error {
	//TODO implement me
	panic("implement me")
}

func (f *ConnectionManagerImpl) GetAllFileNamesFromServer() ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (f *ConnectionManagerImpl) RequestFileFromServer(fileName string) (*os.File, error) {
	//TODO implement me
	panic("implement me")
}

func generateHash(srcFile *os.File) (string, error) {

	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string
	//Open a new hash interface to write to
	hash := md5.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, srcFile); err != nil {
		fmt.Printf("errpr in copy: %v\n", err)
		return returnMD5String, err
	}
	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
