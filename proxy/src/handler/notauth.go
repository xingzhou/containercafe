package handler

import (
	"net/http"
	"io"
	"github.com/golang/glog"
)

//Return 401 for non authorized requests
func NotAuthorizedHandler(w http.ResponseWriter, r *http.Request) {
	glog.Errorf("NotAuthorizedHandler triggered, URI=%s, returning error 401", r.RequestURI)
	w.WriteHeader(401)
	io.WriteString(w, "401 not authorized!\n")
}
