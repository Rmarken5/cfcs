hello:
	echo "Hello"

build-server:
	go build -o file-server  ./server/main

run-server:
	go run ./server/main/server.go --directory ./server/dummy --port 8999

build-client:
	go build -o file-client  ./client/main
run-client:
	go run ./client/main/main.go --directory ./client/client/tmp --port 8999
