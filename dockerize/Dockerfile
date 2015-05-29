FROM ubuntu
RUN apt-get update && apt-get install -y golang
ADD src/ /hijack/src/
WORKDIR /hijack
RUN export GOPATH=/hijack && go build -o bin/hijack src/hijack/hijack.go
EXPOSE 8087
CMD ["/hijack/bin/hijack"]