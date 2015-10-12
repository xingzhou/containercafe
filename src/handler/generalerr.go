package handler

import (
	"net/http"
	"io"
	"fmt"
)

//return general status code
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int){
	Log.Printf("ErrorHandler triggered, URI=%s, returning error %d", r.RequestURI, status)
	w.WriteHeader(status)
	s := fmt.Sprintf("%d error encountered while processing request!\n", status)
	io.WriteString(w, s)
}

func ErrorHandlerWithMsg(w http.ResponseWriter, r *http.Request, status int, msg string){
	Log.Printf("ErrorHandler triggered, URI=%s, error=%d, msg=%s", r.RequestURI, status, msg)
	w.WriteHeader(status)
	s := fmt.Sprintf("%d %s !\n", status, msg)
	io.WriteString(w, s)
}
