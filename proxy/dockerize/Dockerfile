FROM ubuntu

# Add Labels to image that shows git build details in metadata
ARG git_commit_id=unknown
ARG git_remote_url=unknown
ARG build_date=unknown
ARG jenkins_build_id=unknown
ARG jenkins_build_number=unknown

LABEL git-commit-id=${git_commit_id}
LABEL git-remote-url=${git_remote_url}
LABEL build-date=${build_date}
LABEL jenkins-build-id=${jenkins_build_id}
LABEL jenkins-build-number=${jenkins_build_number}

#install docker
#RUN apt-get update && apt-get install -y golang
RUN apt-get update && apt-get install -y wget curl
RUN wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz && rm -f go1.5.1.linux-amd64.tar.gz 

ADD src/ /api-proxy/src/
WORKDIR /api-proxy
RUN export PATH=$PATH:/usr/local/go/bin && export GOPATH=/api-proxy && go build -o bin/api-proxy src/api-proxy/api-proxy.go

ENV log_file_path="/tmp/feeds/logstash_api_proxy"
ENV ccsapi_host=""
ENV ccsapi_uri="/v3/admin/"
ENV ccsapi_compute_node_header="X-Compute-Node"
ENV ccsapi_id_header="X-Container-Id"
ENV ccsapi_id_type_header="X-Id-Type"
ENV docker_port="8089"
ENV docker_api_ver="/v1.18"
ENV tls_inbound="true"
ENV tls_outbound="true"
# containers is started as:
# docker run -v dockerize:/opt/tls_certs --name hjproxy -p 8087:8087 hijack
ENV ca_cert_file="/opt/tls_certs/admin-certs/ca.pem"
ENV client_cert_file="/opt/tls_certs/admin-certs/hjclient.pem"
ENV client_cert_file="/opt/tls_certs/admin-certs/hjclient.pem"
ENV client_key_file="/opt/tls_certs/admin-certs/hjclient.key"
ENV server_cert_file="/opt/tls_certs/admin-certs/hjserver.pem"
ENV server_key_file="/opt/tls_certs/admin-certs/hjserver.key"
ENV kube_admin_key_file="/opt/tls_certs/admin-certs/kadmin.key"
ENV kube_admin_cert_file="/opt/tls_certs/admin-certs/kadmin.pem"
ENV stub_auth_file="/api-proxy/creds.json"
ENV ccsapi_admin_user=""
ENV ccsapi_admin_psswd=""
ENV surrogate_ids="false"
ENV max_container_conn="4"
ENV max_node_conn="0"
ENV kube_port="443"
ENV swarm_node_port="2375"
ENV registry_location=""
ENV registry_admin_password=
ENV consul_ip=""
ENV consul_port=""

# copy all the scripts needed by proxy
ADD make_TLS_certs.sh /api-proxy/make_TLS_certs.sh
ADD mk_user_cert.sh /api-proxy/mk_user_cert.sh
ADD mk_kubeconfig.sh /api-proxy/mk_kubeconfig.sh 
ADD creds.json /api-proxy/creds.json

ADD build.info /api-proxy/build.info

EXPOSE 8087
CMD ["/api-proxy/bin/api-proxy", "8087"]