GOCMD=go
WORKDIR=bin
WORKCONF=$(WORKDIR)/config
JAGCONF=config/jenkins_automation.json
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=jag


all: test build
build:
	mkdir -p $(WORKDIR)
	mkdir -p $(WORKCONF)
	cp $(JAGCONF) $(WORKCONF)
	$(GOBUILD) -o $(WORKDIR)/$(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
