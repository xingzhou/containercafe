FROM ubuntu:16.04

ENV GO_VERSION 1.6.2
ENV K8S_BRANCH merge-28
ENV GO_URL=https://storage.googleapis.com/golang/go$GO_VERSION.linux-amd64.tar.gz
ENV GOPATH=/work
ENV PROJECT_DIR=$GOPATH/src/github.com/kubernetes/kubernetes
ENV KUBERNETES_CONTRIB=mesos
ENV GOROOT=/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$GOROOT/bin:$GOBIN:$PATH
WORKDIR $PROJECT_DIR

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    curl \
    git  \
    rsync \
    build-essential &&\
    mkdir -p $GOPATH/src && mkdir -p $GOPATH/pkg && mkdir -p $GOPATH/bin &&\
    cd / && curl -O -L $GO_URL &&\
    tar zxvf $(basename $GO_URL) &&\
    rm $(basename $GO_URL) &&\
    git clone https://github.com/ibm-contribs/kubernetes.git $PROJECT_DIR &&\
    cd $PROJECT_DIR && git checkout $K8S_BRANCH &&\
    hack/build-go.sh \
    contrib/mesos/cmd/k8sm-scheduler \
    contrib/mesos/cmd/k8sm-executor \
    contrib/mesos/cmd/k8sm-controller-manager \
    contrib/mesos/cmd/km \
    cmd/kubelet \
    cmd/kubectl &&\
    strip -s _output/local/go/bin/km && \
    strip -s _output/local/go/bin/kubelet && \
    strip -s _output/local/go/bin/kubectl && \
    mkdir -p /opt/kubernetes/bin &&\
    cp _output/local/go/bin/km /opt/kubernetes/bin/km &&\
    cp _output/local/go/bin/kubelet /opt/kubernetes/bin/kubelet &&\
    cp _output/local/go/bin/kubectl /opt/kubernetes/bin/kubectl &&\
    rm -rf $PROJECT_DIR &&\
    rm -rf $GOROOT &&\
    apt-get purge --auto-remove -y build-essential rsync git curl &&\
    apt-get autoremove -y &&\
    apt-get clean &&\
    rm -rf /usr/lib/gcc

WORKDIR /opt/kubernetes/bin
