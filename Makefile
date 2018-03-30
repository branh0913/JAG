GOCMD=/usr/local/bin/go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=/usr/local/bin/jag


all: build test
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
build-linux: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
#docker-build: docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v