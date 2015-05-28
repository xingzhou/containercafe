package auth

import (
	"net/http"
	"fmt"
	"strings"
	"httphelper"  // my httphelper
)

var ccsapi_host = "127.0.0.1:8081"					//TODO read from config
var ccsapi_uri = "/v3/containers/auth"				//TODO read from config
var ccsapi_compute_node_header = "X-Compute-Node"	//TODO read from config
var ccsapi_id_header = "X-Container-Id"				//TODO read from config

var Default_redirect_host = "localhost:5000"		//TODO remove this testing default

func get_id_from_uri(uri string, pattern string) string{
	slice1 := strings.SplitAfter(uri, pattern)
	slice2 := strings.Split(slice1[0], "/")
	return slice2[0]
}

func Auth(r *http.Request) (bool, string) {
	ok:=false
	node:=""

	//parse r.RequestURI for container id or exec id
	uri := r.RequestURI
	//1st: look for /containers/exec/<id>/
	id := get_id_from_uri(uri, "/containers/exec/")
	if id == "" {
		//2nd: look for /containers/<id>/
		id = get_id_from_uri(uri, "/containers/")
	}
	if id == "" {
		fmt.Printf("@ Auth: id not found in uri\n")
		//TODO fail here, for now allow a request uri not including <id> to be authenticated
		//fmt.Printf("@ Auth result: %b, %s\n", ok, node)
		//return ok, node
	}

	//forward r header only without body to ccsapi auth endpoint, add X-Container-Id header
	req, _ := http.NewRequest(r.Method, "http://"+ccsapi_host+ccsapi_uri, nil)
	req.Header = r.Header
	req.URL.Host = ccsapi_host
	req.Header.Add(ccsapi_id_header, id)
	client := &http.Client{
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if (err != nil) {
		fmt.Printf("@ Auth: Error in auth request... %v\n", err)

		fmt.Printf("@ Auth result: %b, %s\n", ok, node)
		return ok, node
	}

	//get auth response status, and X-Compute-Node header
	if resp.StatusCode == 200 {
		ok = true
		node = httphelper.GetHeader(resp.Header, ccsapi_compute_node_header)
		if node == ""{
			node = Default_redirect_host   //TODO remove this testing default
		}
	}

	fmt.Printf("@ Auth result: %b, %s\n", ok, node)
	return ok, node
}
