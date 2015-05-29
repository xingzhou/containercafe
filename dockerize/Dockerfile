RUN apt-get update && apt-get install -y golang
ADD src /hijack/src
WORKINGDIR /hijack
RUN export GOPATH=/hijack && go build -o bin/hijack src/hijack.go
EXPOSE 8087
CMD ["/hijack/bin/hijack"]