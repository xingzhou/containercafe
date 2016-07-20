package main

import (
	"github.ibm.com/alchemy-containers/remoteabac/server"
)

func main() {
	server := server.New()
	server.Run()
}
