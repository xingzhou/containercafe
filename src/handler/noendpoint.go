package handler

import (
	"net/http"
	"log"
)

//Return 404 for non-supported URIs
func NoEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("NoEndpointHandler triggered, URI=%s, returning error 404\n", r.RequestURI)
	//w.WriteHeader(404)
	http.NotFound(w, r)
}
