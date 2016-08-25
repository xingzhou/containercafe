package main

import (
	"github.ibm.com/alchemy-containers/openradiant/remoteabac/server"
)

func main() {
	server := server.New()
	server.Run()
}
