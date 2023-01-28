FROM golang:alpine3.16

RUN mkdir "dummy"

WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod .

# download Go modules and dependencies
RUN go mod download

COPY common ./common
COPY server ./server


RUN go mod tidy

RUN ls

# compile application
RUN go build -o fs ./server/main

RUN ls

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8999

ENTRYPOINT [ "./fs",  "--directory", "../dummy", "--port", "8999"]