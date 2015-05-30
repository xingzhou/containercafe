FROM ubuntu
RUN apt-get update && apt-get install -y golang
ADD src/ /hijack/src/
WORKDIR /hijack
RUN export GOPATH=/hijack && go build -o bin/hijack src/hijack/hijack.go

ENV ccsapi_host="127.0.0.1:8081"
ENV ccsapi_uri="/v3/admin/GetComputeNode"
ENV ccsapi_compute_node_header="X-Compute-Node"
ENV ccsapi_id_header="X-Container-Id"
ENV ccsapi_id_type_header="X-Id-Type"
EXPOSE 8087
CMD ["/hijack/bin/hijack", "8087"]