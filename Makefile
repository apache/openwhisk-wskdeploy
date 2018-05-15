#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

SOURCEDIR=.

SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARY=wskdeploy

VERSION=latest

BUILD=`git rev-parse HEAD`

deps:
	@echo "Installing dependencies"
	godep restore -v

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=`git rev-parse HEAD` "

test: deps
	@echo "Testing"
	go test ./... -tags=unit

# Build the project
build:
	go build ${LDFLAGS} -o ${BINARY}

# Run the integration test against OpenWhisk
integration_test:
	@echo "Launch the integration tests."
	go test -v ./... -tags=integration

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

.PHONY: clean install build deps updatedeps format lint test integration_test
