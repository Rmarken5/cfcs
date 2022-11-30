package observer

import "github.com/rmarken5/cfcs/common"

//go:generate mockgen -destination=./mock_observer_test.go -package=observer . Observer
//go:generate mockgen -destination=./mock_subscriber_test.go -package=observer . Subscriber
//go:generate mockgen -destination=./mock_subject_test.go -package=observer . Subject

type Observer interface {
	GetIdentifier() string
	LoadAllFiles(files []common.FileInfo) error
	AddFile(files common.FileInfo) error
}

type Subscriber interface {
	Subscribe(observer Observer)
	Unsubscribe(key string)
	NotifyAllWithFiles(files []common.FileInfo)
	NotifyAllWithFile(file common.FileInfo)
}

type Subject interface {
	AddFile(fileName common.FileInfo)
	RemoveFile(fileName common.FileInfo)
	Subscribe(observer Observer)
	Unsubscribe(key string)
	NotifyAllWithFiles(files []common.FileInfo)
	SetFiles(files []common.FileInfo)
	GetFiles() []common.FileInfo
}
