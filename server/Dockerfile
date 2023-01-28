FROM golang:alpine3.16

WORKDIR /app

COPY . .

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . .

# compile application
RUN go build -o fs ./main

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8999

RUN ls

ENTRYPOINT [ "/fs,  "--directory", "./dummy", "--port", "8999"]