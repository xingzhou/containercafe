# override TLS_OUTBOUND to support K8s
export  tls_outbound=true
# dev-mon01:
#export ccsapi_host=10.140.155.139:8081
export ccsapi_host=10.140.34.174:8081
export consul_ip=10.140.34.174
#export ccsapi_host=10.140.34.164:8081
#export ccsapi_host=localhost:8081
# radiant
# export ccsapi_host=10.140.155.229:8081
# start proxy
bin/hijack 
