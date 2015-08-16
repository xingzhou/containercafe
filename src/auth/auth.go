package auth

import (
	"net/http"
	"log"
	"strings"
	"io/ioutil"
	"encoding/json"
	"httphelper"  	// my httphelper
	"conf"  		// my conf package
)

//returns auth=true/false, compute node name, container/exec id, container id
func Auth(r *http.Request) (ok bool, node string, docker_id string, container string) {

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
		log.Printf("Auth: id not found in uri\n")
		//fail here
		log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
		return ok, node, docker_id, container
	}else{
		//id found in uri
		log.Printf("Auth: id=%s, id_type=%s\n", id, id_type)
	}

	//if type is Exec then get container id from redis cache to forward to getHost api
	var container_id string
	if id_type == "Exec" {
		container_id = conf.RedisGet(id)
	}else{
		container_id = id
	}

	//forward r header only without body to ccsapi auth endpoint, add X-Container-Id header
	req, _ := http.NewRequest("GET", "http://"+conf.GetCcsapiHost()+conf.GetCcsapiUri()+"getHost/"+container_id, nil)
	httphelper.CopyHeader(req.Header, r.Header)  //req.Header = r.Header
	req.URL.Host = conf.GetCcsapiHost()
	req.Header.Add(conf.GetCcsapiIdHeader(), container_id)
	req.Header.Add(conf.GetCcsapiIdTypeHeader(), "Container" /*id_type*/)
	client := &http.Client{
		CheckRedirect: nil,
	}
	resp, err := client.Do(req)
	if (err != nil) {
		log.Printf("Auth: Error in auth request... %v\n", err)

		log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
		return ok, node, docker_id, container
	}

	//get auth response status, and X-Compute-Node header
	if resp.StatusCode == 200 {
		ok = true

		//first check in header
		//node = httphelper.GetHeader(resp.Header, conf.GetCcsapiComputeNodeHeader())
		//second check for json response in body
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)            //Default_redirect_host   //testing default
		if e == nil {
			//convert byte array to string
			//node=string(body[:len(body)])
			log.Printf("Auth: ccsapi raw response=%s\n", body)
			node, container = parse_getHost_Response(body)
		}else {
			//error reading ccsapi response
			log.Printf("Auth result: ok=%t, node='%s'\n", ok, node)
			return ok, node, docker_id, container
		}
		node = node+":"+conf.GetDockerPort()
		//container id needs nova- prefix
		//exec id does not need a prefix
		if id_type == "Container" {
			docker_id = "nova-" + container   //append nova to the id returned from getHost
		}else{//id_type == "Exec"
			docker_id = id  //use the exec id that came in the original req
		}
	}else {
		//demo: exec authentication even if status!=200
		//if id_type == "Exec" {
		//	ok = true
		//	node = conf.Default_redirect_host
		//	docker_id = id
		//}
	}

	log.Printf("Auth result: ok=%t, node=%s, docker_id=%s, container=%s\n", ok, node, docker_id, container)
	return ok, node, docker_id, container
}

//Convert /v*/containers/id/*  to  /<docker_api_ver>/containers/<redirect_resource_id>/*
//Convert /v*/exec/id/*  to  /<docker_api_ver>/exec/<redirect_resource_id>/*
func RewriteURI(reqURI string, redirect_resource_id string) string{
	sl := strings.Split(reqURI, "/")
	redirectURI := "/" + conf.GetDockerApiVer() + "/" + sl[2] + "/" + redirect_resource_id + "/" + sl[4]
	log.Printf("RewriteURI: '%s' --> '%s'\n", reqURI, redirectURI)
	return redirectURI
}

func get_id_from_uri(uri string, pattern string) string{
	var id string
	slice1 := strings.Split(uri, pattern)
	log.Printf("get_id_from_uri: pattern=%s, slice1=%v\n", pattern, slice1)
	if len(slice1) > 1 {
		slice2 := strings.Split(slice1[1], "/")
		id=slice2[0]
	}else{
		id=""
	}
	log.Printf("get_id_from_uri: id=%s\n", id)
	return id
}

func parse_getHost_Response(body []byte) (string, string){

	type Resp struct {
		Container_id  	string
		Container_name 	string
		Host 			string
	}
	var resp Resp

	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Println("parse_getHost_Response: error=%v", err)
	}
	log.Printf("parse_getHost_Response: host=%s, container_id=%s\n", resp.Host, resp.Container_id)
	return resp.Host, resp.Container_id
}
