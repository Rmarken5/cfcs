package observer

/*

func TestFileBroadcastSubject_AddFiles(t *testing.T) {

	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []common.FileInfo{},
		Observers: map[string]Observer{},
		fileMutex: sync.Mutex{},

	}

	fileBroadcastSubject.AddFile("filename.txt")
	fileName := fileBroadcastSubject.Files[0]
	assert.EqualValues(t, "filename.txt", fileName)

}
func TestFileBroadcastSubject_AddFilesFileExists(t *testing.T) {

	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{"filename.txt"},
		Observers: map[string]Observer{},
	}

	fileBroadcastSubject.AddFile("filename.txt")
	fileName := fileBroadcastSubject.Files[0]
	assert.EqualValues(t, "filename.txt", fileName)
	assert.EqualValues(t, len(fileBroadcastSubject.Files), 1)
}

func TestFileBroadcastSubject_RemoveFiles(t *testing.T) {
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{"file1", "file2"},
		Observers: map[string]Observer{},
	}
	fileBroadcastSubject.RemoveFile("file2")
	fileName := fileBroadcastSubject.Files[0]
	assert.EqualValues(t, "file1", fileName)
	assert.EqualValues(t, len(fileBroadcastSubject.Files), 1)
}

func TestFileBroadcastSubject_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	mockObs.EXPECT().GetIdentifier().AnyTimes().Return("obs1")
	mockObs.EXPECT().OnUpdate(gomock.Any())
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{},
		Observers: map[string]Observer{},
	}

	fileBroadcastSubject.Subscribe(mockObs)

	assert.Equal(t, fileBroadcastSubject.Observers["obs1"], mockObs)

}


func TestFileBroadcastSubject_Unsubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{},
		Observers: map[string]Observer{"obs1":mockObs},
	}

	fileBroadcastSubject.Unsubscribe("obs1")

	assert.Len(t, fileBroadcastSubject.Observers, 0)

}

func TestFileBroadcastSubject_NotifyAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	mockObs.EXPECT().OnUpdate([]string{"file1"}).Times(1)
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{"file1"},
		Observers: map[string]Observer{"obs1":mockObs},
	}

	fileBroadcastSubject.NotifyAll()
}
func TestFileBroadcastSubject_NotifyAllUnsub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	mockObs.EXPECT().GetIdentifier().AnyTimes().Return("obs1")
	mockObs.EXPECT().OnUpdate( []string{"file1"}).Return(errors.New("cannot write"))
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{"file1"},
		Observers: map[string]Observer{"obs1":mockObs},
	}
	fileBroadcastSubject.NotifyAll()
	//Test that observer was unsubscribed from
	assert.Len(t, fileBroadcastSubject.Observers, 0)
}

func TestFileBroadcastSubject_GetSetFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockObs := NewMockObserver(ctrl)
	fileBroadcastSubject := FileBroadcastSubject{
		Files:     []string{},
		Observers: map[string]Observer{"obs1":mockObs},
	}
	fileBroadcastSubject.SetFiles([]string{"file1"})
	assert.Contains(t, fileBroadcastSubject.GetFiles(), "file1")
}

*/
