package file_server

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/rmarken5/cfcs/common"
	file_listener "github.com/rmarken5/cfcs/server/file-listener"
	"github.com/rmarken5/cfcs/server/observer"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestServer_AcceptClients(t *testing.T) {
	fmt.Println("hello")
	f, err := os.CreateTemp("", "hello.txt")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	t.Parallel()
	testCases := map[string]struct {
		clientConnectionCommands []string
		mockConn                 func(ctrl *gomock.Controller, address string) net.Conn
		wantErr                  bool
		wantSubscriber           bool
		wantFile                 bool
		mockServer               func(ctrl *gomock.Controller) *Server
	}{
		"should subscribe connection": {
			clientConnectionCommands: []string{fmt.Sprintf("%d\n", common.FILE_LISTENER_CONN_TYPE)},
			mockConn: func(ctrl *gomock.Controller, address string) net.Conn {
				conn, err := net.Dial("tcp", address)
				assert.NoError(t, err, "error getting connection")
				return conn
			},
			mockServer: func(ctrl *gomock.Controller) *Server {
				mockListener := file_listener.NewMockFileListener(ctrl)
				mockSubject := observer.NewMockSubject(ctrl)
				mockSubject.EXPECT().Subscribe(gomock.Any())

				server := Server{
					FileListener: mockListener,
					FileSubject:  mockSubject,
				}

				return &server
			},
			wantErr:        false,
			wantFile:       false,
			wantSubscriber: true,
		},
		"should serve file": {
			clientConnectionCommands: []string{fmt.Sprintf("%d\n", common.FILE_REQUEST_CONN_TYPE), strings.Split(f.Name(), "/")[2]+"\n"},
			mockConn: func(ctrl *gomock.Controller, address string) net.Conn {
				conn, err := net.Dial("tcp", address)
				assert.NoError(t, err, "error getting connection")
				return conn
			},
			mockServer: func(ctrl *gomock.Controller) *Server {
				mockListener := file_listener.NewMockFileListener(ctrl)
				mockSubject := observer.NewMockSubject(ctrl)

				server := Server{
					FileListener: mockListener,
					FileSubject:  mockSubject,
					FileDirectory: "/"+strings.Split(f.Name(), "/")[1],
				}

				return &server
			},
			wantErr:        false,
			wantFile:       false,
			wantSubscriber: true,
		},
		"should return reading buffer": {
			clientConnectionCommands: []string{""},
			mockConn: func(ctrl *gomock.Controller, address string) net.Conn {
				conn, err := net.Dial("tcp", address)
				assert.NoError(t, err, "error getting connection")
				return conn
			},
			mockServer: func(ctrl *gomock.Controller) *Server {
				mockListener := file_listener.NewMockFileListener(ctrl)
				mockSubject := observer.NewMockSubject(ctrl)

				server := Server{
					FileListener: mockListener,
					FileSubject:  mockSubject,
				}

				return &server
			},
			wantErr: true,
		},
		"should return reading value not integer": {
			clientConnectionCommands: []string{"hello\n"},
			mockConn: func(ctrl *gomock.Controller, address string) net.Conn {
				conn, err := net.Dial("tcp", address)
				assert.NoError(t, err, "error getting connection")
				return conn
			},
			mockServer: func(ctrl *gomock.Controller) *Server {
				mockListener := file_listener.NewMockFileListener(ctrl)
				mockSubject := observer.NewMockSubject(ctrl)

				server := Server{
					FileListener: mockListener,
					FileSubject:  mockSubject,
				}

				return &server
			},
			wantErr: true,
		},
		"should return reading value - violating protocol": {
			clientConnectionCommands: []string{"99\n"},
			mockConn: func(ctrl *gomock.Controller, address string) net.Conn {
				conn, err := net.Dial("tcp", address)
				assert.NoError(t, err, "error getting connection")
				return conn
			},
			mockServer: func(ctrl *gomock.Controller) *Server {
				mockListener := file_listener.NewMockFileListener(ctrl)
				mockSubject := observer.NewMockSubject(ctrl)

				server := Server{
					FileListener: mockListener,
					FileSubject:  mockSubject,
				}

				return &server
			},
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			a, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
			if err != nil {
				fmt.Println(err)
				return
			}
			l, err := net.ListenTCP("tcp", a)
			if err != nil {
				fmt.Println(err)
				return
			}
			server := tc.mockServer(ctrl)
			go server.AcceptClients(l)

			conn := tc.mockConn(ctrl, l.Addr().String())

			for _, command := range tc.clientConnectionCommands {
				_, err := conn.Write([]byte(command))
				time.Sleep(time.Second) //Wait for the go routine to read connection
				assert.NoError(t, err, "error writing to connection %s", command)
			}
			time.Sleep(time.Second) //Wait for the go routine to subscribe connection

		})
	}

}

/*
func Test_ListenForFiles(t *testing.T) {

}

func Test_AddFilesToSubject(t *testing.T) {

}
*/
