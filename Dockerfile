#
# Dockerfile for building containers-ng-ansible image
#
FROM ubuntu:14.04

MAINTAINER Paolo Dettori <dettori@us.ibm.com>

# image metadata
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

WORKDIR /containers-ng-ansible
ENV ANSIBLE_LIBRARY /containers-ng-ansible/sl-vms-raw/library
ENV ANSIBLE_INVENTORY /containers-ng-ansible/sl-vms-raw/library/sl_inventory.py
ENV SL_INVENTORY /containers-ng-ansible/sl-vms-raw/library/sl_inventory.py

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -yq \
     vim \
     python-pip \
     python-dev \
     build-essential \
     libffi-dev \
     libssl-dev \
     ssh \
     python-openssl

RUN python -m pip install ansible==1.9.5 \
                softlayer==4.1.1 \
#                pyopenssl \
                ndg-httpsclient \
                pyasn1

# config ssh
RUN mkdir ~/.ssh && echo "Host * " > ~/.ssh/config && \
    echo "   StrictHostKeyChecking no" >> ~/.ssh/config && \
    echo "   UserKnownHostsFile /dev/null" >> ~/.ssh/config && \
    echo "   LogLevel ERROR" >> ~/.ssh/config

COPY . /containers-ng-ansible
