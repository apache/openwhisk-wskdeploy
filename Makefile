SOURCEDIR=.

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARY=wskdeploy

VERSION=1.0.0

BUILD=`git rev-parse HEAD`

deps:
	@echo "Installing dependencies"
	@go get -d -t ./...


LDFLAGS=-ldflags "-X main.Version=`date -u '+%Y-%m-%dT%H:%M:%S'` -X main.Build=`git rev-parse HEAD` "

updatedeps:
	@echo "Updating all dependencies"
	@go get -d -u -f -fix -t ./...


test: deps
	@echo "Testing"
	go test ./...

# Build the project
build: test
	go build ${LDFLAGS} -o ${BINARY}

format:
	@echo "Formatting"
	go fmt ./...

lint: format
	@echo "Linting"
	golint .

install:
	go install

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY}; fi

.PHONY: clean install