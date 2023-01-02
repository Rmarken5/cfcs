package file_listener

import (
	"github.com/fsnotify/fsnotify"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"testing"
)

func TestFileListener_ListenForFiles(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		assert.NoError(t, err)
	}
	defer watcher.Close()
	fileListener := FileListener{
		watcher,
	}

	event, err := fileListener.ListenForFiles("")
	assert.NoError(t, err)
	assert.NotNil(t, event)
}

func TestFileListener_ReadDirectory(t *testing.T) {
	ctr := gomock.NewController(t)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	entry := NewMockDirEntry(ctr)
	entry.EXPECT().IsDir().AnyTimes().Return(false)
	entry.EXPECT().Name().AnyTimes().Return("Derp")

	fileListener := FileListener{
		watcher,
	}
	dirEntries := []fs.DirEntry{
		entry,
	}
	files := fileListener.ReadDirectory(dirEntries)

	assert.Equal(t, files[0], "Derp")
}

func TestFileListener_ReadDirectoryIsDir(t *testing.T) {
	ctr := gomock.NewController(t)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()
	entry := NewMockDirEntry(ctr)
	entry.EXPECT().IsDir().AnyTimes().Return(true)
	entry.EXPECT().Name().AnyTimes().Return("Derp")

	fileListener := FileListener{
		watcher,
	}
	dirEntries := []fs.DirEntry{
		entry,
	}
	files := fileListener.ReadDirectory(dirEntries)

	assert.Len(t, files, 0)
}

func TestFileListener_BuildFileInfoFromPath(t *testing.T) {
	wantDir := "temp-test"
	err := os.Mkdir(wantDir, 0777)
	if err != nil {
		assert.NoError(t, err, "should not get error creating dir")
	}
	defer os.RemoveAll(wantDir)
	testCases := map[string]struct {
		wantError    bool
		wantFileName string
	}{
		"should get successful file info from file path": {
			wantError:    false,
			wantFileName: "tmp.txt",
		},
	}

	for name, tt := range testCases {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			err2 := os.WriteFile(wantDir+"/"+tt.wantFileName, []byte(""), 0666)
			assert.NoError(t, err2)
			info, err := BuildFileInfoFromPath(wantDir + "/" + tt.wantFileName)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, info)
				assert.Equal(t, tt.wantFileName, info.FileName)
			}
		})
	}
}
