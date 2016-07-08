# override TLS_OUTBOUND to support K8s
export tls_outbound=true
export tls_inbound=true
export server_key_file=certs/hjserver-key.pem 
export server_cert_file=certs/hjserver.pem
# dev-mon01:
export ccsapi_host=10.140.34.174:8081
export consul_ip=10.140.34.174
# radiant
# start proxy
bin/hijack 
