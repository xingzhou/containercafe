#!/bin/bash

docker run -d --name etcd -p 4001:4001 quay.io/coreos/etcd:v2.0.12 --listen-client-urls http://0.0.0.0:4001 --advertise-client-urls http://0.0.0.0:4001
sleep 2

docker exec etcd /etcdctl set /abac-policy '
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "scheduler", "nonResourcePath": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "scheduler", "namespace": "*", "resource": "*", "apiGroup": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "controller-manager", "nonResourcePath": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "controller-manager", "namespace": "*", "resource": "*", "apiGroup": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "kubelet", "nonResourcePath": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "kubelet", "namespace": "*", "resource": "*", "apiGroup": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "admin", "nonResourcePath": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "admin", "namespace": "*", "resource": "*", "apiGroup": "*"}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "/api", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "/api/*/", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "/apis/", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "/apis/*", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "/apis/*/*", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "/apis/*/*", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "*", "nonResourcePath": "version", "readonly": true}}
{"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": {"user": "alice", "namespace": "alice", "resource": "*", "apiGroup": "*"}}
'

# On a Mac if docker is started via Virtualbox, use `docker-machine env` to find out the ip address of the VM and use that as ETCD_IP
# On a Linux machine if docker is started natively, use 127.0.0.1
#ETCD_IP=192.168.99.100
ETCD_IP=127.0.0.1

# standalone mode
remoteabac --address=:8888 --tls-cert-file=apiserver.pem --tls-private-key-file=apiserver-key.pem --authorization-policy-file=etcd@http://${ETCD_IP}:4001/abac-policy

# containerized mode
#docker run -p 8888:8888 -v `pwd`:/tmp haih/remoteabac --address=:8888 --tls-cert-file=/tmp/apiserver.pem --tls-private-key-file=/tmp/apiserver-key.pem --authorization-policy-file=etcd@http://localhost:4001/abac-policy
