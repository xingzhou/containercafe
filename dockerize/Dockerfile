FROM ubuntu

#install docker
#RUN apt-get update && apt-get install -y golang
RUN apt-get update && apt-get install -y wget
RUN wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz && rm -f go1.5.1.linux-amd64.tar.gz 

ADD src/ /hijack/src/
WORKDIR /hijack
RUN export PATH=$PATH:/usr/local/go/bin && export GOPATH=/hijack && go build -o bin/hijack src/hijack/hijack.go

ENV log_file_path="/tmp/feeds/logstash_hijack_proxy"
ENV ccsapi_host="127.0.0.1:8081"
ENV ccsapi_uri="/v3/admin/"
ENV ccsapi_compute_node_header="X-Compute-Node"
ENV ccsapi_id_header="X-Container-Id"
ENV ccsapi_id_type_header="X-Id-Type"
ENV docker_port="8089"
ENV docker_api_ver="/v1.18"
ENV tls_inbound="true"
ENV tls_outbound="true"
ENV client_cert_file="/opt/tls_certs/hjclient.pem"
ENV client_key_file="/opt/tls_certs/hjclient.key"
ENV server_cert_file="/opt/tls_certs/hjserver.pem"
ENV server_key_file="/opt/tls_certs/hjserver.key"
ENV stub_auth_file="/opt/tls_certs/creds.json"
ENV kube_admin_key_file = "/opt/tls_certs/kadmin.key"
ENV kube_admin_cert_file = "/opt/tls_certs/kadmin.pem"
ENV ccsapi_admin_user=""
ENV ccsapi_admin_psswd=""
ENV surrogate_ids="false"
ENV max_container_conn="4"
ENV max_node_conn="0"
ENV kube_port="443"
ENV swarm_node_port="2375"
#ENV registry_auth_tok="eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiIyMzAxODkiLCJhdXRoIjoiIiwiZW1haWwiOiJzd2FybUBkZXYudGVzdCIsInNlcnZlcmFkZHJlc3MiOiJyZWdpc3RyeS1pY2UtZGV2LXRlc3Quc3RhZ2UxLm5nLmJsdWVtaXgubmV0In0="
ENV registry_location="registry-ice-dev-test.stage1.ng.bluemix.net"
ENV registry_admin_password=230189
ENV consul_ip="127.0.0.1"
ENV consul_port="8500"

#Copy the self-signed certs and keys
ADD hjserver.pem /hijack/hjserver.pem
ADD hjserver.key  /hijack/hjserver.key
ADD hjclient.pem /hijack/hjclient.pem
ADD hjclient.key  /hijack/hjclient.key
ADD kadmin.key /hijack/kadmin.key
ADD kadmin.pem /hijack/kadmin.pem

ADD build.info /hijack/build.info

EXPOSE 8087
CMD ["/hijack/bin/hijack", "8087"]
