FROM ubuntu
RUN apt-get update && apt-get install -y golang
ADD src/ /hijack/src/
WORKDIR /hijack
RUN export GOPATH=/hijack && go build -o bin/hijack src/hijack/hijack.go

ENV ccsapi_host="127.0.0.1:8081"
ENV ccsapi_uri="/v3/admin/"
ENV ccsapi_compute_node_header="X-Compute-Node"
ENV ccsapi_id_header="X-Container-Id"
ENV ccsapi_id_type_header="X-Id-Type"
ENV docker_port="8089"
ENV docker_api_ver="v1.17"
ENV tls_inbound="true"
ENV tls_outbound="true"
ENV cert_file="cert.pem"
ENV key_file="key.pem"
ENV ccsapi_admin_user ""
ENV ccsapi_admin_psswd ""

#Copy the self-signed server cert and key
ADD cert.pem /hijack/cert.pem
ADD key.pem  /hijack/key.pem

EXPOSE 8087
CMD ["/hijack/bin/hijack", "8087"]