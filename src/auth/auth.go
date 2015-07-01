package auth

import (
	"net/http"
	"fmt"
	"strings"
	"httphelper"  // my httphelper
	"os"
	"io/ioutil"
	"encoding/json"
)

//TODO move env settings to separate conf module
//ENV defaults
var ccsapi_host = "127.0.0.1:8081"
var ccsapi_uri = "/v3/admin/getHost/"   //ends with id
var ccsapi_compute_node_header = "X-Compute-Node"
var ccsapi_id_header = "X-Container-Id"
var ccsapi_id_type_header = "X-Id-Type"				//Container or Exec
var docker_port = "8089"
var docker_api_ver = "v1.17"
var tls_inbound = "false"
var tls_outbound = "false"
var cert_file = "cert.pem"
var key_file = "key.pem"

var default_listen_port = 8087
var Default_redirect_host = "vizio-dev-host2:8089"		//TODO remove this testing default

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
	load_env_var("docker_api_ver", &docker_api_ver)
	load_env_var("tls_inbound", &tls_inbound)
	load_env_var("tls_outbound", &tls_outbound)
	load_env_var("cert_file", &cert_file)
	load_env_var("key_file", &key_file)
}

func GetDefaultListenPort() int {
	return default_listen_port
}

func SetTlsInbound (b bool){
	if b {
		tls_inbound = "true"
	}else {
		tls_inbound = "false"
	}
}

func IsTlsInbound() bool {
	if tls_inbound == "true" {
		return true
	} else {
		return false
	}
}

func GetCertFile() string{
	return cert_file
}

func GetKeyFile() string{
	return key_file
}

//returns auth=true/false, compute node name, container/exec id
func Auth(r *http.Request) (bool, string, string) {
	ok:=false
	node:=""
	docker_id:=""

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
		fmt.Printf("@ Auth result: ok=%t, node='%s'\n", ok, node)
		return ok, node, docker_id
	}else{
		fmt.Printf("@ Auth: id=%s, id_type=%s\n", id, id_type)
	}

	//forward r header only without body to ccsapi auth endpoint, add X-Container-Id header
	req, _ := http.NewRequest(r.Method, "http://"+ccsapi_host+ccsapi_uri+id, nil)
	httphelper.CopyHeader(req.Header, r.Header)  //req.Header = r.Header
	req.URL.Host = ccsapi_host
	req.Header.Add(ccsapi_id_header, id)
	req.Header.Add(ccsapi_id_type_header, id_type)
	client := &http.Client{
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if (err != nil) {
		fmt.Printf("@ Auth: Error in auth request... %v\n", err)

		fmt.Printf("@ Auth result: ok=%t, node='%s'\n", ok, node)
		return ok, node, docker_id
	}

	//get auth response status, and X-Compute-Node header
	if resp.StatusCode == 200 {
		ok = true
		//first check in header
		node = httphelper.GetHeader(resp.Header, ccsapi_compute_node_header)
		if node == "" {
			//second check for json response in body
			defer resp.Body.Close()
			body, e := ioutil.ReadAll(resp.Body)            //Default_redirect_host   //testing default
			if e == nil {
				//convert byte array to string
				//node=string(body[:len(body)])
				fmt.Printf("@ Auth: ccsapi raw response=%s\n", body)
					//fmt.Printf("@ Auth: ccsapi response=%s", httphelper.PrettyJson(body))
				node, docker_id = parse_getHost_Response(body)
			}else {
				//error reading ccsapi response
				fmt.Printf("@ Auth result: ok=%t, node='%s'\n", ok, node)
				return ok, node, docker_id
			}
		}
		node = node+":"+docker_port
		if id_type == "Container" {
			//container id needs nova- prefix
			//exec id does not need a prefix
			docker_id = "nova-" + docker_id
		}
	}else {
		//TODO remove the following demo exec authentication even if status!=200
		if id_type == "Exec" {
			ok = true
			node = Default_redirect_host
			docker_id = id
		}
	}

	fmt.Printf("@ Auth result: ok=%t, node='%s', docker_id='%s'\n", ok, node, docker_id)
	return ok, node, docker_id
}

//Convert /v*/containers/id/*  to  /<docker_api_ver>/containers/<redirect_resource_id>/*
//Convert /v*/exec/id/*  to  /<docker_api_ver>/exec/<redirect_resource_id>/*
func RewriteURI(reqURI string, redirect_resource_id string) string{
	sl := strings.Split(reqURI, "/")
	redirectURI := "/" + docker_api_ver + "/" + sl[2] + "/" + redirect_resource_id + "/" + sl[4]
	fmt.Printf("@ RewriteURI: '%s' --> '%s'\n", reqURI, redirectURI)
	return redirectURI
}

func get_id_from_uri(uri string, pattern string) string{
	var id string
	slice1 := strings.Split(uri, pattern)
	fmt.Printf("@ get_id_from_uri: pattern=%s, slice1=%v\n", pattern, slice1)
	if len(slice1) > 1 {
		slice2 := strings.Split(slice1[1], "/")
		id=slice2[0]
	}else{
		id=""
	}
	fmt.Printf("@ get_id_from_uri: id=%s\n", id)
	return id
}

func parse_getHost_Response(body []byte) (string, string){

	type Resp struct {
		Container_id  	string
		Container_name 	string
		Host 			string
	}
	var resp Resp

	fmt.Printf("@ parse_getHost_Response: json=%s\n", body)
	err := json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("@ parse_getHost_Response: error=%v", err)
	}
	fmt.Printf("@ parse_getHost_Response: host=%s, container_id=%s\n", resp.Host, resp.Container_id)
	return resp.Host, resp.Container_id
}
