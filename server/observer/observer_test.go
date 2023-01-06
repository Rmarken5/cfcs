package observer

import (
	"github.com/golang/mock/gomock"
	"github.com/rmarken5/cfcs/common"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestConnectionObserver_LoadAllFiles(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		wantError        bool
		mockDependencies func(controller *gomock.Controller, files []common.FileInfo) *ConnectionObserver
		mockFileInfos    []common.FileInfo
	}{
		"successfully load all files": {
			wantError: false,
			mockDependencies: func(controller *gomock.Controller, files []common.FileInfo) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(len(SERVER_SENDING_FILE_LIST.String())+1, nil)
				c.EXPECT().Write(gomock.Any()).Times(len(files)).Return(1, nil)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
			mockFileInfos: []common.FileInfo{{
				FileName: "file1",
				Hash:     "hash1",
			}},
		},
		"fails to initiate communication": {
			wantError: true,
			mockDependencies: func(controller *gomock.Controller, files []common.FileInfo) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(0, net.ErrClosed)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
			mockFileInfos: nil,
		},
		"fails to load files": {
			wantError: true,
			mockDependencies: func(controller *gomock.Controller, files []common.FileInfo) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(len(SERVER_SENDING_FILE_LIST.String())+1, nil)
				c.EXPECT().Write(gomock.Any()).Return(0, net.ErrWriteToConnected)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
			mockFileInfos: []common.FileInfo{{
				FileName: "file1",
				Hash:     "hash1",
			}},
		},
		"fails on second file write": {
			wantError: true,
			mockDependencies: func(controller *gomock.Controller, files []common.FileInfo) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(len(SERVER_SENDING_FILE_LIST.String())+1, nil)
				c.EXPECT().Write(gomock.Any()).Return(1, nil)
				c.EXPECT().Write(gomock.Any()).Return(0, net.ErrWriteToConnected)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
			mockFileInfos: []common.FileInfo{{
				FileName: "file1",
				Hash:     "hash1",
			},
				{
					FileName: "file2",
					Hash:     "hash2",
				}},
		},
	}

	for name, tt := range testCases {
		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			connObs := tt.mockDependencies(ctrl, tt.mockFileInfos)
			err := connObs.LoadAllFiles(tt.mockFileInfos)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConnectionObserver_AddFile(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		wantError        bool
		mockDependencies func(controller *gomock.Controller) *ConnectionObserver
	}{
		"successfully add file": {
			wantError: false,
			mockDependencies: func(controller *gomock.Controller) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(1, nil)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
		},
		"fails to initiate communication": {
			wantError: true,
			mockDependencies: func(controller *gomock.Controller) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(0, net.ErrClosed)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
		},
		"fails to load files": {
			wantError: true,
			mockDependencies: func(controller *gomock.Controller) *ConnectionObserver {
				c := NewMockConn(controller)
				c.EXPECT().Write(gomock.Any()).Return(0, net.ErrWriteToConnected)

				return &ConnectionObserver{
					Address: "TestAddress",
					Conn:    c,
				}
			},
		},
	}

	for name, tt := range testCases {
		name := name
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			connObs := tt.mockDependencies(ctrl)
			err := connObs.AddFile(common.FileInfo{})

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConnectionData_GetIdentifier(t *testing.T) {
	connData := ConnectionObserver{
		"hello",
		nil,
	}

	addr := connData.GetIdentifier()
	assert.EqualValues(t, "hello", addr)
}
