package handler

import (
	"net/http"
	"io"
)

//Return 501 for non-supported URIs
func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	Log.Printf("NotImplementedHandler triggered, URI=%s, returning error 501", r.RequestURI)
	w.WriteHeader(501)
	io.WriteString(w,"Not implemented\n")
}
