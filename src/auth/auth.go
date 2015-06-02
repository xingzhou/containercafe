package auth

import (
	"net/http"
	"fmt"
	"strings"
	"httphelper"  // my httphelper
	"os"
	"io/ioutil"
)

var ccsapi_host = "127.0.0.1:8081"
var ccsapi_uri = "/v3/admin/getHost/"   //ends with id
var ccsapi_compute_node_header = "X-Compute-Node"
var ccsapi_id_header = "X-Container-Id"
var ccsapi_id_type_header = "X-Id-Type"				//Container or Exec
var docker_port="5000"

var Default_redirect_host = "localhost:5000"		//TODO remove this testing default

func load_env_var(env_name string, target *string) {
	s:=os.Getenv(env_name)
	if s != "" {
		*target = s
	}
	fmt.Printf("@ load_env_var: %s=%s\n",env_name, *target)
}

func LoadEnv(){
	load_env_var("ccsapi_host", &ccsapi_host)
	load_env_var("ccsapi_uri", &ccsapi_uri)
	load_env_var("ccsapi_compute_node_header", &ccsapi_compute_node_header)
	load_env_var("ccsapi_id_header", &ccsapi_id_header)
	load_env_var("ccsapi_id_type_header", &ccsapi_id_type_header)
	load_env_var("docker_port", &docker_port)
}

func get_id_from_uri(uri string, pattern string) string{
	slice1 := strings.Split(uri, pattern)
	slice2 := strings.Split(slice1[1], "/")
	return slice2[0]
}

func Auth(r *http.Request) (bool, string) {
	ok:=false
	node:=""

	//parse r.RequestURI for container id or exec id
	uri := r.RequestURI
	//1st: look for /containers/<id>/
	id := get_id_from_uri(uri, "/containers/")
	id_type:="Container"
	if id == "" {
		//2nd: look for /exec/<id>/
		id = get_id_from_uri(uri, "/exec/")
		id_type="Exec"
	}
	if id == "" {
		fmt.Printf("@ Auth: id not found in uri\n")
		//fail here, for now allow a request uri not including <id> to be authenticated
		fmt.Printf("@ Auth result: %b, node='%s'\n", ok, node)
		return ok, node
	}else{
		fmt.Printf("@ Auth: id=%s, id_type=%s\n", id, id_type)
	}

	//forward r header only without body to ccsapi auth endpoint, add X-Container-Id header
	req, _ := http.NewRequest(r.Method, "http://"+ccsapi_host+ccsapi_uri+id, nil)
	req.Header = r.Header
	req.URL.Host = ccsapi_host
	req.Header.Add(ccsapi_id_header, id)
	req.Header.Add(ccsapi_id_type_header, id_type)
	client := &http.Client{
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if (err != nil) {
		fmt.Printf("@ Auth: Error in auth request... %v\n", err)

		fmt.Printf("@ Auth result: %b, node='%s'\n", ok, node)
		return ok, node
	}

	//get auth response status, and X-Compute-Node header
	if resp.StatusCode == 200 {
		ok = true
		//first check for header
		node = httphelper.GetHeader(resp.Header, ccsapi_compute_node_header)
		if node == ""{
			//second check for text response
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)			//Default_redirect_host   //testing default
			//TODO err check
			//convert byte array to string
			node=string(body[:len(body)])
			node= node + ":" + docker_port
		}
	}

	fmt.Printf("@ Auth result: %b, node='%s'\n", ok, node)
	return ok, node
}
