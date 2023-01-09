package main

import (
	"fmt"
	"github.com/golang/mock/gomock"
	file_listener "github.com/rmarken5/cfcs/server/file-listener"
	"github.com/rmarken5/cfcs/server/observer"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestServer_AcceptClients(t *testing.T) {
	fmt.Println("hello")
	t.Parallel()
	testCases := map[string]struct{
		clientConnectionCommands []string
		wantErr bool
		wantSubscriber bool
		wantFile bool
	}{
		"should subscribe connection" : {
			clientConnectionCommands: []string{fmt.Sprintf("%d",observer.FILE_LISTENER_CONN_TYPE)},
			wantErr: false,
			wantFile: false,
			wantSubscriber: true,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mockListener := file_listener.NewMockFileListener(ctrl)
			mockSubject := observer.NewMockSubject(ctrl)
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

			server := server{
				FileListener: mockListener,
				FileSubject:  mockSubject,
			}
			server.acceptClients(l)

			conn, err := net.Dial("tcp", a.String())
			assert.NoError(t, err, "error getting connection")

			defer func(conn net.Conn) {
				err := conn.Close()
				assert.NoError(t, err, "error closing connection")
			}(conn)


			if tc.wantSubscriber {
				for _, command := range tc.clientConnectionCommands {
					n, err := conn.Write([]byte(command + "\n"))
					assert.NoError(t, err, "error writing to connection %s", command)
					assert.Greater(t, n, 0)
				}
			}
		})
	}

}
/*
func Test_ListenForFiles(t *testing.T) {

}

func Test_AddFilesToSubject(t *testing.T) {

}
*/
