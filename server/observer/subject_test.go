package observer

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rmarken5/cfcs/common"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestFileBroadcastSubject_AddFiles(t *testing.T) {

	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{},
		Observers: map[string]Observer{},
		fileMutex: sync.RWMutex{},
	}

	fileInfo := common.FileInfo{
		FileName: "filename.txt",
		Hash:     "",
	}

	fileBroadcastSubject.AddFile(fileInfo)
	fileName := fileBroadcastSubject.Files[0].FileName
	assert.EqualValues(t, "filename.txt", fileName)

}
func TestFileBroadcastSubject_AddFilesFileExists(t *testing.T) {

	fileInfo := common.FileInfo{
		FileName: "filename.txt",
		Hash:     "",
	}

	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{fileInfo},
		Observers: map[string]Observer{},
	}

	fileBroadcastSubject.AddFile(fileInfo)
	fileName := fileBroadcastSubject.Files[0].FileName
	assert.EqualValues(t, "filename.txt", fileName)
	assert.EqualValues(t, len(fileBroadcastSubject.Files), 1)
}

func TestFileBroadcastSubject_RemoveFiles(t *testing.T) {
	fileInfoOne := common.FileInfo{
		FileName: "filename1.txt",
		Hash:     "1",
	}
	fileInfoTwo := common.FileInfo{
		FileName: "filename2.txt",
		Hash:     "2",
	}

	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{fileInfoOne, fileInfoTwo},
		Observers: map[string]Observer{},
	}
	fileBroadcastSubject.RemoveFile(fileInfoTwo)
	fileName := fileBroadcastSubject.Files[0].FileName
	assert.EqualValues(t, "filename1.txt", fileName)
	assert.EqualValues(t, len(fileBroadcastSubject.Files), 1)
}

func TestFileBroadcastSubject_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := map[string]struct {
		mockFileBroadcastSubject func(ctrl *gomock.Controller) *FileBroadcastSubject
	}{
		"successful": {
			mockFileBroadcastSubject: func(ctrl *gomock.Controller) *FileBroadcastSubject {
				mockObs := NewMockObserver(ctrl)
				mockObs.EXPECT().GetIdentifier().AnyTimes().Return("obs1")
				mockObs.EXPECT().LoadAllFiles(gomock.Any()).Return(nil)
				fileBroadcastSubject := FileBroadcastSubject{
					Files: []common.FileInfo{
						{
							FileName: "filename1.txt",
							Hash:     "1",
						},
					},
					Observers: map[string]Observer{"obs1": mockObs},
				}
				return &fileBroadcastSubject
			},
		},
		"error loading files": {
			mockFileBroadcastSubject: func(ctrl *gomock.Controller) *FileBroadcastSubject {
				mockObs := NewMockObserver(ctrl)
				mockObs.EXPECT().GetIdentifier().AnyTimes().Return("obs1")
				mockObs.EXPECT().LoadAllFiles(gomock.Any()).Return(errors.New("cannot load files"))
				fileBroadcastSubject := FileBroadcastSubject{
					Files: []common.FileInfo{
						{
							FileName: "filename1.txt",
							Hash:     "1",
						},
					},
					Observers: map[string]Observer{"obs1": mockObs},
				}
				return &fileBroadcastSubject
			},
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			sub := tc.mockFileBroadcastSubject(ctrl)
			sub.Subscribe(sub.Observers["obs1"])
		})
	}
}

func TestFileBroadcastSubject_Unsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{},
		Observers: map[string]Observer{"obs1": mockObs},
	}

	fileBroadcastSubject.Unsubscribe("obs1")

	assert.Len(t, fileBroadcastSubject.Observers, 0)

}

func TestFileBroadcastSubject_NotifyAllWithFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	mockObs.EXPECT().LoadAllFiles(gomock.Any()).Return(nil)
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{{FileName: "file1"}},
		Observers: map[string]Observer{"obs1": mockObs},
	}

	fileBroadcastSubject.NotifyAllWithFiles([]common.FileInfo{{FileName: "file2"}})
}

func TestFileBroadcastSubject_GetSetFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{},
		Observers: map[string]Observer{"obs1": mockObs},
	}
	fileBroadcastSubject.SetFiles([]common.FileInfo{{FileName: "file1"}})
	assert.Contains(t, fileBroadcastSubject.GetFiles(), common.FileInfo{FileName: "file1"})
}
