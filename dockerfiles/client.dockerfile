FROM golang:alpine3.16

RUN mkdir "files"

WORKDIR /cfcs

# copy Go modules and dependencies to image
COPY go.mod .
COPY common ./common
COPY client ./client

# download Go modules and dependencies
RUN go mod tidy

# compile application
RUN go build -o c ./client/main

# tells Docker that the container listens on specified network ports at runtime
ARG DEFAULT_HOST="localhost"
ENV SERVER_HOST=$DEFAULT_HOST

ENTRYPOINT ./c --directory ../files --host $SERVER_HOST --port 8999
