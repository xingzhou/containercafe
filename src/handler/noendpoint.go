package handler

import (
	"net/http"
	"log"
	"io"
)

//Return 404 for non-supported URIs
func NoEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("NoEndpointHandler triggered, URI=%s, returning error 404\n", r.RequestURI)
	//http.NotFound(w, r)
	w.WriteHeader(404)
	io.WriteString(w, "404 not found (Hijackproxy NoEndpointHandler)\n")
}
