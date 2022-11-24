package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/golang/mock/gomock"
	"testing"
)

/*func Test_AcceptClients(t *testing.T) {

	ctr := gomock.NewController(t)
	netConn := NewMockConn(ctr)
	addr := NewMockAddr(ctr)
	netConn.EXPECT().RemoteAddr().AnyTimes().Return(addr)
	netListener := NewMockListener(ctr)
	netListener.EXPECT().Close().Return(nil)
	netListener.EXPECT().Accept().AnyTimes().Return(netConn,nil)
	fileListener := NewMockIFileListener(ctr)
	subject := NewMockSubject(ctr)
	server := server{
		FileListener: fileListener,
		FileSubject: subject,
	}

	server.acceptClients(netListener)
}*/

func Test_EvaluateEvent(t *testing.T) {
	create := fsnotify.Event{
		Name: "/file1",
		Op:   1,
	}
	remove := fsnotify.Event{
		Name: "/file1",
		Op:   4,
	}
	listenerChannel := make(chan fsnotify.Event)
	ctr := gomock.NewController(t)
	fileListener := NewMockIFileListener(ctr)
	subject := NewMockSubject(ctr)
	subject.EXPECT().AddFile("file1").Times(1)
	subject.EXPECT().RemoveFile("file1").Times(1)

	server := server{
		FileListener: fileListener,
		FileSubject:  subject,
	}

	go server.evaluateEvent(listenerChannel)
	listenerChannel <- create
	go server.evaluateEvent(listenerChannel)
	listenerChannel <- remove

}
