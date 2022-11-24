package observer

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectionData_OnUpdate_Nil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockConn(ctrl)
	m.EXPECT().Write(gomock.Any()).Return(11, nil)

	conData := ConnectionData{
		"addr",
		m,
	}

	err := conData.LoadAllFiles([]string{"file-1.txt"})
	if err != nil {
		t.FailNow()
	}
}

func TestConnectionData_OnUpdate_NotNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockConn(ctrl)
	m.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("SomeError"))

	conData := ConnectionData{
		"addr",
		m,
	}

	err := conData.LoadAllFiles([]string{"file-1.txt"})
	if err == nil {
		t.FailNow()
	}
}

func TestConnectionData_GetIdentifier(t *testing.T) {
	connData := ConnectionData{
		"hello",
		nil,
	}

	addr := connData.GetIdentifier()
	assert.EqualValues(t, "hello", addr)
}
