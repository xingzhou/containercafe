FROM ubuntu

# Add Labels to image that shows git build details in metadata
ARG git_commit_id=unknown
ARG git_remote_url=unknown
ARG build_date=unknown
ARG jenkins_build_id=unknown
ARG jenkins_build_number=unknown

LABEL git-commit-id=${git_commit_id} \
      git-remote-url=${git_remote_url} \
      build-date=${build_date} \
      jenkins-build-id=${jenkins_build_id} \
      jenkins-build-number=${jenkins_build_number}

#install docker
#RUN apt-get update && apt-get install -y golang
RUN apt-get update && apt-get install -y wget curl
RUN wget --quiet https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz && rm -f go1.5.1.linux-amd64.tar.gz 

ADD src/ /api-proxy/src/
WORKDIR /api-proxy
RUN export PATH=$PATH:/usr/local/go/bin && export GOPATH=/api-proxy && go build -o bin/api-proxy src/api-proxy/api-proxy.go

ENV log_file_path="/tmp/feeds/logstash_api_proxy" \
    docker_port="8089" \
    docker_api_ver="/v1.18" \
    tls_inbound="true" \
    tls_outbound="true" \
    ca_cert_file="/opt/tls_certs/admin-certs/ca.pem" \
    client_cert_file="/opt/tls_certs/admin-certs/hjclient.pem" \
    client_cert_file="/opt/tls_certs/admin-certs/hjclient.pem" \
    client_key_file="/opt/tls_certs/admin-certs/hjclient.key" \
    server_cert_file="/opt/tls_certs/admin-certs/hjserver.pem" \
    server_key_file="/opt/tls_certs/admin-certs/hjserver.key" \
    kube_admin_key_file="/opt/tls_certs/admin-certs/kadmin.key" \
    kube_admin_cert_file="/opt/tls_certs/admin-certs/kadmin.pem" \
    stub_auth_file="/opt/tls_certs/creds.json" \
    surrogate_ids="false" \
    max_container_conn="4" \
    max_node_conn="0" \
    kube_port="443" \
    swarm_node_port="2375" \
    registry_location="" \
    registry_admin_password= \
    consul_ip="" \
    consul_port=""

# copy all the scripts needed by proxy
ADD make_TLS_certs.sh /api-proxy/make_TLS_certs.sh
ADD mk_user_cert.sh /api-proxy/mk_user_cert.sh
ADD mk_kubeconfig.sh /api-proxy/mk_kubeconfig.sh 

ADD build.info /api-proxy/build.info

EXPOSE 8087
CMD ["/api-proxy/bin/api-proxy", "8087"]