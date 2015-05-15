GOPATH := $(CURDIR)
all: build

build:
	GOPATH=$(GOPATH) go fmt src/command/client/client.go
	GOPATH=$(GOPATH) go fmt src/command/server/server.go
	GOPATH=$(GOPATH) go build -o bin/client src/command/client/client.go
	GOPATH=$(GOPATH) go build -o bin/server src/command/server/server.go

run:
	go fmt src/TimeServer.go
	GOPATH=$(GOPATH) go run src/command/client/client.go
	GOPATH=$(GOPATH) go run src/command/server/server.go
 
install:
	GOPATH=$(GOPATH) go install command/client
	GOPATH=$(GOPATH) go install command/server

client:
	GOPATH=$(GOPATH) go build -o bin/client src/command/client/client.go

server:
	GOPATH=$(GOPATH) go build -o bin/server src/command/server/server.go

rclient:
	./bin/client

rserver:
	./bin/server