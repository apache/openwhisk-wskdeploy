FROM golang:1.8

# Install zip
RUN apt-get -y update && \
    apt-get -y install zip emacs

ENV GOPATH=/

# Download and install tools
RUN echo "Installing the godep tool"
RUN go get github.com/tools/godep

ADD . /src/github.com/apache/incubator-openwhisk-wskdeploy

# Load all of the dependencies from the previously generated/saved godep generated godeps.json file
RUN echo "Restoring Go dependencies"
RUN cd /src/github.com/apache/incubator-openwhisk-wskdeploy && /bin/godep restore -v

# All of the Go CLI binaries will be placed under a build folder
RUN rm -rf /src/github.com/apache/incubator-openwhisk-wskdeploy/build
RUN mkdir /src/github.com/apache/incubator-openwhisk-wskdeploy/build

ARG WSKDEPLOY_OS
ARG WSKDEPLOY_ARCH

# Build the Go wsk CLI binaries and compress resultant binaries
RUN chmod +x /src/github.com/apache/incubator-openwhisk-wskdeploy/build.sh
RUN cd /src/github.com/apache/incubator-openwhisk-wskdeploy/ && ./build.sh
