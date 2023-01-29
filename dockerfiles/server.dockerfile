FROM golang:alpine3.16

RUN mkdir "dummy"

WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod .
COPY common ./common
COPY server ./server

RUN go mod tidy

# compile application
RUN go build -o fs ./server/main

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8999

ENTRYPOINT [ "./fs",  "--directory", "../dummy", "--port", "8999"]