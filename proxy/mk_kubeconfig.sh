#!/bin/bash

PATH="$1"
TENANT_ID="$2"

cd $PATH


echo "apiVersion: v1" > kube-config
echo "clusters:" >> kube-config
echo "- cluster:" >> kube-config
echo "    server: https://localhost:8087" >> kube-config
echo "    certificate-authority: ca.pem" >> kube-config
echo "  name: radiant" >> kube-config
echo "contexts:" >> kube-config
echo "- context:" >> kube-config
echo "    cluster: radiant" >> kube-config
echo "    namespace: s$TENANT_ID-default" >> kube-config
echo "    user: $TENANT_ID" >> kube-config
echo "  name: radiant" >> kube-config
echo "current-context: radiant" >> kube-config
echo "kind: Config" >> kube-config
echo "preferences: {}" >> kube-config
echo "users:" >> kube-config
echo "- name: $TENANT_ID" >> kube-config
echo "  user:" >> kube-config
echo "    client-certificate: cert.pem" >> kube-config
echo "    client-key: key.pem" >> kube-config