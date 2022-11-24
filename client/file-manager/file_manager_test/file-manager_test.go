package file_manager_test

import (
	"git.mills.io/prologic/bitcask"
	file_manager "github.com/rmarken5/cfcs/client/file-manager"

	"testing"
)

/*func TestDownloadFile(t *testing.T) {

	ctr := gomock.NewController(t)
	mockConn := NewMockConn(ctr)


	fileManager := file_manager.FileManagerImpl{}
	err := fileManager.DownloadFile(mockConn, "/dev/null")
	if err != nil {
		t.Errorf("TestDownloadFile: %v\n", err)
		t.Fail()
	}
}*/

func TestWriteFileNameToDatabase(t *testing.T) {
	db, err := bitcask.Open("./test-db")
	defer func(db *bitcask.Bitcask) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	if err != nil {
		t.Fail()
	}
	fileManager := file_manager.NewFileManagerImpl([]string{}, db)
	err = fileManager.WriteFileHashToDB("file-1")

	if err != nil {
		t.Fail()
	}

}

func TestWriteFileToQueue(t *testing.T) {
	var queue []string
	manager := file_manager.NewFileManagerImpl(queue, nil)
	manager.AddFileToQueue("file1")
}

func TestRemoveFileFromQueue(t *testing.T) {
	queue := []string{"file1"}
	manager := file_manager.NewFileManagerImpl(queue, nil)
	manager.RemoveFileFromQueue("file1")
}
