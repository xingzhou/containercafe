#!/bin/bash

NUM="$1"
TENANT_ID="$2"
FILE_PATH="$3"

if [ ! -f "$FILE_PATH" ]; then
    touch "$FILE_PATH"
    echo "apiVersion: v1" >> "$FILE_PATH"
    echo "kind: Pod" >> "$FILE_PATH"
    echo "metadata:" >> "$FILE_PATH"
    echo "  name: $TENANT_ID-kube-web-test$NUM" >> "$FILE_PATH"
    echo "  labels:" >> "$FILE_PATH"
    echo "    app: web-ms-demo" >> "$FILE_PATH"
    echo "  annotations:" >> "$FILE_PATH"
    echo "    aaa: bbbb" >> "$FILE_PATH"
    echo "    containers-label.alpha.kubernetes.io/com.swarm.tenant.0: sf7f413cb-a678-412d-b024-8e17e28bcb88-default" >> "$FILE_PATH"
    echo "spec:" >> "$FILE_PATH"
    echo "  containers:" >> "$FILE_PATH"
    echo "    - name: kube-web-server" >> "$FILE_PATH"
    echo "      image: mrsabath/web-ms:v3" >> "$FILE_PATH"
    echo "      ports:" >> "$FILE_PATH"
    echo "        - containerPort: 80" >> "$FILE_PATH"
    echo "      env:" >> "$FILE_PATH"
    echo "        -" >> "$FILE_PATH"
    echo "         name: \"TEST\"" >> "$FILE_PATH"
    echo "         value: \"web-test$NUM\"" >> "$FILE_PATH"
fi