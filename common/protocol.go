package common

type ConnHandlerMessage int

const (
	FILE_LISTENER_CONN_TYPE ConnHandlerMessage = iota
	FILE_REQUEST_CONN_TYPE
	SERVER_READY_TO_RECEIVE_FILE_REQUEST
	SERVER_SENDING_FILE_LIST
)

func IsCHM(n int) bool {
	conv := ConnHandlerMessage(n)
	return conv.String() != ""
}

func (chm ConnHandlerMessage) String() string {
	switch chm {
	case FILE_LISTENER_CONN_TYPE:
		return "FILE_LISTENER_CONNECTION"
	case FILE_REQUEST_CONN_TYPE:
		return "FILE_REQUEST_CONNECTION"
	case SERVER_READY_TO_RECEIVE_FILE_REQUEST:
		return "SERVER_READY_TO_RECEIVE_FILE_REQUEST"
	}
	return ""
}
