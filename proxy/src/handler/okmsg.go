package handler

import (
	"net/http"
	"io"
	"fmt"
)

//return general status code
func OkHandler(w http.ResponseWriter, r *http.Request, status int){
	Log.Printf("OkHandler triggered, URI=%s, returning status %d", r.RequestURI, status)
	w.WriteHeader(status)
	//s := fmt.Sprintf("%d success while processing request!\n", status)
	//io.WriteString(w, s)
}

func OkHandlerWithMsg(w http.ResponseWriter, r *http.Request, status int, msg string){
	Log.Printf("OkHandlerWithMsg triggered, URI=%s, status=%d, msg=%s", r.RequestURI, status, msg)
	w.WriteHeader(status)
	s := fmt.Sprintf("%d %s!\n", status, msg)
	io.WriteString(w, s)
}
