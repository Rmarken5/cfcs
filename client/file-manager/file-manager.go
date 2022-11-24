package file_manager

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"strings"
)

type FileManagerImpl struct {
	fileQueue map[string]bool
	db        *bitcask.Bitcask
	sftp      *sftp.Client
}

func NewFileManagerImpl(dbLocation string, sftp *sftp.Client) *FileManagerImpl {
	open, err := bitcask.Open(dbLocation, bitcask.WithSync(true), bitcask.WithAutoRecovery(true))
	if err != nil {
		log.Fatalf("Error opening db: %v\n", err)
	}

	return &FileManagerImpl{
		fileQueue: make(map[string]bool),
		db:        open,
		sftp:      sftp,
	}
}

// WriteSrcToDest takes a file and writes it to the given fullyQualifiedPath of the directory.
func (f *FileManagerImpl) WriteSrcToDest(srcFile *sftp.File, fullyQualifiedPath string) (string, error) {
	// TODO: Make fullyQualifiedPath a directory.
	// TODO: Make dstPath come from a configuration as fullyQualifiedPath that's passed into function rather than hardcoded.
	dstPath := "/home/ryan/programming/go-programs/file-client/tmp"
	parts := strings.Split(fullyQualifiedPath, "/")
	fileName := parts[len(parts)-1]

	// TODO: Use srcFile and fullyQualified path to build dstFile
	dstFile, err := os.Create(dstPath + "/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// Copy the file
	_, err = srcFile.WriteTo(dstFile)
	return fileName, err
}

func (f *FileManagerImpl) GetSourceFile(sshConfig *ssh.ClientConfig, fullyQualifiedPath string) (*sftp.File, error) {

	// Open the source file
	srcFile, err := f.sftp.Open(fullyQualifiedPath)
	return srcFile, err
}

func (f *FileManagerImpl) WriteFileHashToDB(fileName string, srcFile *sftp.File) error {
	hash, err := generateHash(*srcFile)
	if err != nil {
		return err
	}
	err = f.db.Put([]byte(fileName), []byte(hash))
	fmt.Println(f.db.Get([]byte(fileName)))
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
	f.sftp.Close()
}

func (f *FileManagerImpl) ShouldWriteToDB(fileName string, srcFile *sftp.File) bool {

	fmt.Println("filename lookup: " + fileName)
	if !f.db.Has([]byte(fileName)) {
		return true
	}
	fmt.Println("db has file")

	hash, err := generateHash(*srcFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return false
	}
	fmt.Println("Hash: " + hash)
	hashToCompare, err := f.db.Get([]byte(fileName))
	fmt.Println("File hash: " + string(hashToCompare))
	return hash != string(hashToCompare)
}
func generateHash(srcFile sftp.File) (string, error) {

	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string
	//Open a new hash interface to write to
	hash := md5.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, &srcFile); err != nil {
		fmt.Printf("errpr in copy: %v\n", err)
		return returnMD5String, err
	}
	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
