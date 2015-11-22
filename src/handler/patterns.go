package handler

import (
	"strings"
	"net/http"

	"auth"
)

//return true if uri pattern is supported
func IsSupportedPattern(uri string, patterns []string) bool{
	for i:=0; i < len(patterns); i++ {
		//if uri contains patterns[i]
		if strings.Contains(uri, patterns[i]) {
			return true
		}
	}
	return false
}

//return uri prefix pattern
func GetUriPattern(uri string, patterns []string) string{
	for i:=0; i < len(patterns); i++ {
		if strings.Contains(uri, patterns[i]) {
			return patterns[i]
		}
	}
	return ""
}


// TODO: New advanced patterns with simple expressions

type RouteHandler func(w http.ResponseWriter, r *http.Request, body []byte, creds auth.Creds, vars map[string]string, req_id string)

type Route struct{
	method 		string
	pattern 	string
	handler		RouteHandler
}

func newRoute(method string, pattern string, handler RouteHandler) Route{
	return Route{method, pattern, handler}
}

func TestPatt() {
	Log.Print("TestPatt start >>>>>")

	//dummy router 1
	var route Route
	route.method = "GET"
	route.pattern = "/{version}/containers/{id}/json"
	route.handler = containers_json
	var(
		w http.ResponseWriter
		r *http.Request
		body []byte
		creds auth.Creds
		vars map[string]string
		req_id string = "123"
	)
	route.handler(w, r, body, creds, vars, req_id)

	//dummy router 2
	var routes []Route
	routes = []Route{
		newRoute("GET", "/{version}/containers/{id}/json", containers_json),
		newRoute("DELETE", "/{version}/containers/{id}/json", containers_json),
	}
	routes[0].handler(w, r, body, creds, vars, "321")

	Log.Print("TestPatt end <<<<<")
}

func containers_json(w http.ResponseWriter, r *http.Request, body []byte, creds auth.Creds, vars map[string]string, req_id string){
	Log.Print("containers_json start >>>>>")
	Log.Printf("containers_json req_id=%s", req_id)
	Log.Print("containers_json end <<<<<")
}

