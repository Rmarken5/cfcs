FROM golang:alpine3.16

WORKDIR /app

COPY . .

RUN go build -o file-client  ./main

ENTRYPOINT [ "./file-client" ]

