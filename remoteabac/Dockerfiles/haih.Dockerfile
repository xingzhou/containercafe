FROM scratch

MAINTAINER Hai Huang <haih@us.ibm.com>

COPY remoteabac /opt/kubernetes/
COPY ruser /opt/kubernetes/
COPY empty /tmp/
ENTRYPOINT ["/opt/kubernetes/remoteabac"]

