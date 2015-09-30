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
	s := fmt.Sprintf("%d error encountered while processing request (Hijackproxy ErrorHandler)\n", status)
	io.WriteString(w, s)
}

func ErrorHandlerWithMsg(w http.ResponseWriter, r *http.Request, status int, msg string){
	log.Printf("ErrorHandler triggered, URI=%s, error=%d, msg=%s\n", r.RequestURI, status, msg)
	w.WriteHeader(status)
	s := fmt.Sprintf("%d %s (Hijackproxy ErrorHandler)\n", status, msg)
	io.WriteString(w, s)
}
