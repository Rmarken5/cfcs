// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rmarken5/cfcs/server/file-listener (interfaces: FileListener)

// Package file_listener is a generated GoMock package.
package file_listener

import (
	fs "io/fs"
	reflect "reflect"

	fsnotify "github.com/fsnotify/fsnotify"
	gomock "github.com/golang/mock/gomock"
)

// MockFileListener is a mock of FileListener interface.
type MockFileListener struct {
	ctrl     *gomock.Controller
	recorder *MockFileListenerMockRecorder
}

// MockFileListenerMockRecorder is the mock recorder for MockFileListener.
type MockFileListenerMockRecorder struct {
	mock *MockFileListener
}

// NewMockFileListener creates a new mock instance.
func NewMockFileListener(ctrl *gomock.Controller) *MockFileListener {
	mock := &MockFileListener{ctrl: ctrl}
	mock.recorder = &MockFileListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileListener) EXPECT() *MockFileListenerMockRecorder {
	return m.recorder
}

// ListenForFiles mocks base method.
func (m *MockFileListener) ListenForFiles(arg0 string) (chan fsnotify.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenForFiles", arg0)
	ret0, _ := ret[0].(chan fsnotify.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListenForFiles indicates an expected call of ListenForFiles.
func (mr *MockFileListenerMockRecorder) ListenForFiles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenForFiles", reflect.TypeOf((*MockFileListener)(nil).ListenForFiles), arg0)
}

// ReadDirectory mocks base method.
func (m *MockFileListener) ReadDirectory(arg0 []fs.DirEntry) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadDirectory", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// ReadDirectory indicates an expected call of ReadDirectory.
func (mr *MockFileListenerMockRecorder) ReadDirectory(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDirectory", reflect.TypeOf((*MockFileListener)(nil).ReadDirectory), arg0)
}