FROM scratch

MAINTAINER Hai Huang <haih@us.ibm.com>

ARG git_commit_id=unknown
ARG build_id=unknown
ARG build_number=unknown
ARG build_date=unknown
ARG git_tag=unknown
ARG git_remote_url=unknown

LABEL git-commit-id=${git_commit_id}
LABEL build-id=${build_id}
LABEL build-number=${build_number}
LABEL build-date=${build_date}
LABEL git-tag=${git_tag}
LABEL git-remote-url=${git_remote_url}

#RUN mkdir -p /opt/kubernetes

COPY remoteabac /opt/kubernetes/
COPY ruser /opt/kubernetes/
COPY empty /tmp/
ENTRYPOINT ["/opt/kubernetes/remoteabac"]
