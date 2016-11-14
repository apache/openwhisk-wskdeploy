FROM golang:1.7

# Install zip
RUN apt-get -y update && \
    apt-get -y install zip

# Download and install tools
RUN echo "Installing the godep tool"
RUN go get github.com/tools/godep


RUN echo "Installing the go-bindata tool"
RUN go get github.com/jteeuwen/go-bindata/...

RUN echo "Installing the openwhisk go-client dependency"
RUN go get github.com/openwhisk/openwhisk-client-go/whisk


RUN cd $GOPATH/src/github.com/openwhisk

RUN git clone http://github.com/openwhisk/wskdeploy $GOPATH/src/github.com/openwhisk/wskdeploy

RUN cd $GOPATH/src/github.com/openwhisk/wskdeploy && git checkout development

RUN cd $GOPATH/src/github.com/openwhisk/wskdeploy && go get -t -d ./...

RUN cd $GOPATH/src/github.com/openwhisk/wskdeploy && make build

RUN cd $GOPATH/src/github.com/openwhisk/wskdeploy && make install