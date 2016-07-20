FROM ubuntu:16.04

ENV GO_VERSION 1.6.2
ENV GO_URL=https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz
ENV GOPATH=/work
ENV PROJECT_DIR=$GOPATH/src/github.ibm.com/alchemy-containers/remoteabac
ENV GOROOT=/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOROOT/bin:$GOBIN:$PATH
WORKDIR $PROJECT_DIR

COPY . $PROJECT_DIR

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    git  \
    curl \
    build-essential &&\
    mkdir -p $GOPATH/src && mkdir -p $GOPATH/pkg && mkdir -p $GOPATH/bin &&\
    cd / && curl -O -L $GO_URL &&\
    tar zxvf $(basename $GO_URL) &&\
    rm $(basename $GO_URL) &&\
    go get -u github.com/tools/godep &&\
    cd $PROJECT_DIR/cmd/remoteabac && godep go install -a remoteabac.go &&\
    cd $PROJECT_DIR/cmd/ruser && godep go install -a ruser.go &&\
    mkdir -p /opt/kubernetes &&\
    cp $GOBIN/remoteabac /opt/kubernetes/ &&\
    cp $GOBIN/ruser /opt/kubernetes/ &&\
    rm -rf $GOPATH &&\
    rm -rf $GOROOT &&\
    apt-get purge --auto-remove -y build-essential rsync git curl &&\
    apt-get autoremove -y &&\
    apt-get clean &&\
    rm -rf /usr/lib/gcc

WORKDIR /opt/kubernetes/
ENTRYPOINT ["/opt/kubernetes/remoteabac"]
