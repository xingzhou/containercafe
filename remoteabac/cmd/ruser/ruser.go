package main

import (
	"log"

	"github.ibm.com/alchemy-containers/openradiant/remoteabac/policy"
)

func main() {
	policy, err := policy.New()
	if err != nil {
		log.Fatalf("Received an error: %v\n", err)
	}
	policy.ProcessRequest()
}
