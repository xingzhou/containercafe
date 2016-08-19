#!/bin/bash

function make_kubeconfig {
    local CONFIG_DIR="$1"
    local TENANT_ID="$2"

    cat > "$CONFIG_DIR/kube-config" <<YAML
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8087
    certificate-authority: ca.pem
  name: radiant
contexts:
- context:
    cluster: radiant
    namespace: s${TENANT_ID}-default
    user: $TENANT_ID
  name: radiant
current-context: radiant
kind: Config
preferences: {}
users:
- name: ${TENANT_ID}
  user:
    client-certificate: cert.pem
    client-key: key.pem
YAML
}

make_kubeconfig "$@"
