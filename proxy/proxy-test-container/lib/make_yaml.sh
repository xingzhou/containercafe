#!/bin/bash

NUM="$1"
TENANT_ID="$2"
FILE_PATH="$3"

if [ ! -f "$FILE_PATH" ]; then
    cat > "$FILE_PATH" <<YAML
apiVersion: v1
kind: Pod
metadata:
  name: $TENANT_ID-kube-web-test$NUM
  labels:
    app: web-ms-demo
  annotations:
    aaa: bbbb
    containers-label.alpha.kubernetes.io/com.swarm.tenant.0: sf7f413cb-a678-412d-b024-8e17e28bcb88-default
spec:
  containers:
    - name: kube-web-server
      image: mrsabath/web-ms:v3
      ports:
        - containerPort: 80
      env:
        -
         name: "TEST"
         value: "web-test$NUM"
YAML
fi
