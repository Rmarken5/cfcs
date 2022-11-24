package observer

//go:generate mockgen -destination=./mock_observer_test.go -package=observer . Observer
//go:generate mockgen -destination=./mock_subscriber_test.go -package=observer . Subscriber
//go:generate mockgen -destination=./mock_subject_test.go -package=observer . Subject

type Observer interface {
	GetIdentifier() string
	LoadAllFiles(files []string) error
	AddFile(files string) error
}

type Subscriber interface {
	Subscribe(Observer)
	Unsubscribe(string)
	NotifyAllWithFiles(files []string)
	NotifyAllWithFile(file string)
}

type Subject interface {
	AddFile(fileName string)
	RemoveFile(fileName string)
	Subscribe(observer Observer)
	Unsubscribe(key string)
	NotifyAllWithFiles(files []string)
	SetFiles(files []string)
	GetFiles() []string
}
