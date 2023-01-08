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
	t.Parallel()
	testCases := map[string]struct {
		wantDirectory string
		wantErr       bool
	}{
		"successfully returns event channel": {
			wantDirectory: "",
			wantErr:       false,
		},
		"Should throw error for directory": {
			wantDirectory: "/this/should/not/work",
			wantErr:       true,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			watcher, err := fsnotify.NewWatcher()
			assert.NoError(t, err)
			defer func() {
				err2 := watcher.Close()
				assert.NoError(t, err2)
			}()
			fileListener := FileListener{
				watcher,
			}

			event, err := fileListener.ListenForFiles(tc.wantDirectory)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, event)
			}
		})
	}
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
	testDir := "temp-test"
	err := os.Mkdir(testDir, 0777)
	if err != nil {
		assert.NoError(t, err, "should not get error creating dir")
	}
	defer func() {
		err2 := os.RemoveAll(testDir)
		assert.NoError(t, err2)
	}()
	testCases := map[string]struct {
		wantError bool
		wantFile  string
		wantDir   string
	}{
		"should get successful file info from file path": {
			wantError: false,
			wantFile:  "tmp.txt",
			wantDir:   testDir,
		},
		"should handle error opening file" :{
			wantError: true,
			wantFile: "tmp.txt",
			wantDir: "should/not/work",
		},

	}

	for name, tt := range testCases {
		name := name
		tt := tt
		t.Run(name, func(t *testing.T) {
			err := os.WriteFile(testDir+"/"+tt.wantFile, []byte(""), 0666)
			assert.NoError(t, err)
			defer func () {
				err := os.Remove(testDir + "/" + tt.wantFile)
				assert.NoError(t, err)
			}()
			info, err := BuildFileInfoFromPath(tt.wantDir + "/" + tt.wantFile)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, info)
				assert.Equal(t, tt.wantFile, info.FileName)
			}
		})
	}
}
