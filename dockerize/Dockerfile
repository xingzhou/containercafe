FROM ubuntu
RUN apt-get update && apt-get install -y golang
ADD src/ /hijack/src/
WORKDIR /hijack
RUN export GOPATH=/hijack && go build -o bin/hijack src/hijack/hijack.go

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
ENV ccsapi_admin_user=""
ENV ccsapi_admin_psswd=""
ENV surrogate_ids="false"
ENV max_container_conn="4"
ENV max_node_conn="0"
ENV kube_port="8080"

#Copy the self-signed certs and keys
ADD hjserver.pem /hijack/hjserver.pem
ADD hjserver.key  /hijack/hjserver.key
ADD hjclient.pem /hijack/hjclient.pem
ADD hjclient.key  /hijack/hjclient.key

ADD build.info /hijack/build.info

EXPOSE 8087
CMD ["/hijack/bin/hijack", "8087"]
