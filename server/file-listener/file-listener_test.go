package file_listener

import (
	"github.com/fsnotify/fsnotify"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"testing"
)


func TestFileListener_ListenForFiles(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatal(err)
	}
	defer watcher.Close()


	fileListener := FileListener {
		watcher,
	}

	event := fileListener.ListenForFiles("")
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

	fileListener := FileListener {
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

	fileListener := FileListener {
		watcher,
	}
	dirEntries := []fs.DirEntry{
		entry,
	}
	files := fileListener.ReadDirectory(dirEntries)

	assert.Len(t, files, 0)
}
