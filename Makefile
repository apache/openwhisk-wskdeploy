SOURCEDIR=.
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARY=wsktool

VERSION=1.0.0

BUILD=`git rev-parse HEAD`

deps:
	@echo "Installing dependencies"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

updatedeps:
	@echo "Updating all dependencies"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u


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
	go install ${LDFLAGS}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY}; fi

.PHONY: clean install