package handler

import (
	"net/http"
	"log"
	"fmt"

	"conf"
)

//Return 404 for all non-supported URIs
func NoEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("NoEndpointHandler triggered, URI=%s, returning error 404\n", r.RequestURI)
	//w.WriteHeader(404)
	http.NotFound(w, r)
}

func HealthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HealthEndpointHandler triggered, URI=%s\n", r.RequestURI)
	fmt.Fprintf(w,"hjproxy up\n")
	fmt.Fprintf(w,"This instance served %d requests\n", conf.GetNumServedRequests())
}
