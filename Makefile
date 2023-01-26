.PHONY: hello
hello:
	echo "Hello"

.PHONY: build-server
build-server:
	go build -o file-server  ./server/main

.PHONY: run-server
run-server:
	go run ./server/main/server.go --directory ./server/dummy --port 8999

.PHONY: build-client
build-client:
	go build -o file-client  ./client/main

.PHONY: run-client
run-client:
	go run ./client/main/main.go --directory ./client/tmp --port 8999

.PHONY: gen
gen:
	go generate -v -x ./...