package handler

import (
	"strings"
	"net/http"

	"auth"
)

//return true if uri prefix is supported
func IsSupportedPattern(uri string, prefixes []string) bool{
	for i:=0; i < len(prefixes); i++ {
		//if uri contains patterns[i]
		if strings.Contains(uri, prefixes[i]) {
			return true
		}
	}
	return false
}

//return uri prefix if supported
func GetUriPattern(uri string, prefixes []string) string{
	for i:=0; i < len(prefixes); i++ {
		if strings.Contains(uri, prefixes[i]) {
			return prefixes[i]
		}
	}
	return ""
}

////////////////////////////////////////////////////
// Routing based on patterns with simple expressions
////////////////////////////////////////////////////

type RouteHandler func(w http.ResponseWriter, r *http.Request, body []byte, creds auth.Creds, vars map[string]string, req_id string)

type Route struct{
	method 		string
	pattern 	string
	handler		RouteHandler
}

type Router struct{
	routes		[]Route
}

func NewRoute(method string, pattern string, handler RouteHandler) Route{
	return Route{method, pattern, handler}
}

func NewRouter(routes []Route) *Router{
	router := new (Router)
	router.routes = routes
	return router
}

func (router *Router) CheckRoute(req *http.Request)(found bool, route Route){
	for _, route = range router.routes {
		if found, _ = route.matchRoute(req); found {
			return
		}
	}
	return
}

//uses SelectRoute to determine target handler and invokes it
func (router *Router) DoRoute(w http.ResponseWriter, req *http.Request, body []byte, creds auth.Creds, req_id string) {
	for  _, route := range router.routes {
		if found, vars := route.matchRoute(req); found {
			route.handler(w, req, body, creds, vars, req_id)
			return
		}
	}
	Log.Printf("No route found req_id=%s", req_id)
}

func (route *Route) matchRoute(req *http.Request) (match bool, vars map[string]string){
	match = false
	if (route.method != req.Method) && (route.method != "*") {
		return
	}

	//see if it is a wildcard pattern
	if route.pattern == "*" {
		match = true
		return
	}

	pattern_parts := strings.Split(route.pattern, "/")
	uri_parts := strings.Split(req.RequestURI, "/")

	//check pattern length vs uri length
	if len(pattern_parts) != len(uri_parts) {
		return
	}

	//match part by part and fill vars as you go
	vars = make(map[string]string)
	for i:=0; i < len(pattern_parts); i++ {
		if strings.Contains(pattern_parts[i], "{") {
			vars[pattern_parts[i]] = uri_parts[i]
			continue
		}
		if pattern_parts[i] != uri_parts[i]{
			return
		}
	}

	match = true
	return
}

////////////////////////////
// Test code
///////////////////////////

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
		NewRoute("GET", "/{version}/containers/{id}/json", containers_json),
		NewRoute("DELETE", "/{version}/containers/{id}/json", containers_json),
	}
	routes[0].handler(w, r, body, creds, vars, "321")

	Log.Print("TestPatt end <<<<<")
}

func containers_json(w http.ResponseWriter, r *http.Request, body []byte, creds auth.Creds, vars map[string]string, req_id string){
	Log.Print("containers_json start >>>>>")
	Log.Printf("containers_json req_id=%s", req_id)
	Log.Print("containers_json end <<<<<")
}

