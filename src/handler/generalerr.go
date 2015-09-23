package handler

import (
	"net/http"
	"log"
	"io"
	"fmt"
)

//return general status code
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int){
	log.Printf("ErrorHandler triggered, URI=%s, returning error %d\n", r.RequestURI, status)
	w.WriteHeader(status)
	s := fmt.Sprintf("%d error encountered will processing request (Hijackproxy ErrorHandler)\n", status)
	io.WriteString(w, s)
}
