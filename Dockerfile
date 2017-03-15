FROM golang:1.7

# Install zip
RUN apt-get -y update && \
    apt-get -y install zip emacs

# Download and install tools
RUN echo "Installing the godep tool"
RUN go get github.com/tools/godep


RUN echo "Installing the go-bindata tool"
RUN go get github.com/jteeuwen/go-bindata/...

RUN echo "Installing the openwhisk go-client dependency"
RUN go get github.com/openwhisk/openwhisk-client-go/whisk 

RUN cd $GOPATH/src/github.com/openwhisk/openwhisk-client-go/whisk && git fetch && git checkout i18n

RUN git clone http://github.com/paulcastro/openwhisk-wskdeploy $GOPATH/src/github.com/openwhisk/openwhisk-wskdeploy

RUN cd $GOPATH/src/github.com/openwhisk/openwhisk-wskdeploy && git checkout issue_160

RUN cd $GOPATH/src/github.com/openwhisk/openwhisk-wskdeploy && go get -t -d ./...

RUN cd $GOPATH/src/github.com/openwhisk/openwhisk-wskdeploy && make build

RUN cd $GOPATH/src/github.com/openwhisk/openwhisk-wskdeploy && make install