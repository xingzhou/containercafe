package handler

import (
	"net/http"
	"log"

	"auth"
	"conf"
)

// handler for Kubernetes
func KubeEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	log.Printf("------> KubeEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// check if uri pattern is accepted
	if ! auth.VerifyKubeUriPattern(r.RequestURI){
		log.Printf("Kube pattern not accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
		log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Call Auth interceptor
	ok, node, namespace := auth.KubeAuth(r)
	if !ok {
		log.Printf("Authentication failed for req_id=%s", req_id)
		log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Handle request
	kubeHandler(w, r, node, namespace, req_id)

	log.Printf("------ Completed processing of request req_id=%s\n", req_id)
}


func kubeHandler(w http.ResponseWriter, r *http.Request, redirect_host string, redirect_resource_id string, req_id string) {
	//TODO
	return
}

