package auth

import (
	"net/http"
	"log"
	"strings"
	"strconv"

	"conf"  		// my conf package
)

// returns auth=true/false, compute node name, container/exec id, container id,
// override tls flag is used in swarm case only
func DockerAuth(r *http.Request) (status int, node string, docker_id string,
	container string, tls_override bool) {

	//parse r.RequestURI for container id or exec id
	uri := r.RequestURI
	id, id_type := get_id_and_type(uri)
	//if type is Exec then get container id from redis cache to forward to getHost api
	var container_id string
	if id_type == "Exec" {
		container_id = conf.RedisGet(id)
	}else{
		container_id = id
	}

	var host GetHostResp
	if id_type == "None" {
		status, host = getHost(r, "NoneContainer")
	}else{
		status, host = getHost(r, container_id)
	}
	if status != 200 {
		log.Printf("Auth result: status=%d\n", status)
		return status, "", "", "", false
	}
	node = host.Host + ":" + conf.GetDockerPort()
	container = host.Container_id
	//container id needs nova- prefix
	//exec id does not need a prefix
	if id_type == "Container" {
		docker_id = "nova-"+container   //append nova to the id returned from getHost
	}
	if id_type == "Exec" {
		docker_id = id  //use the exec id that came in the original req
	}
	if id_type == "None" {
		docker_id = ""
	}
	tls_override = false
	if host.Swarm {
		node = host.Mgr_host    //Mgr_host = host:port
		//insert space_id in the header to be forwarded
		r.Header.Set("X-Auth-Token", host.Space_id)
		if id_type == "Container" {
			docker_id = container
		}
		if id_type == "Exec" {
			docker_id = id  //use the exec id that came in the original req
		}
		if id_type == "None" {
			docker_id = ""
		}
		if !host.Swarm_tls{
			tls_override = true  // no tls for this outbound req regardless of proxy conf
		}
		if host.Host != ""{
			// node swarm node directly if known
			node = host.Host + ":" + strconv.Itoa(conf.GetSwarmNodePort())
		}
	}
	log.Printf("Auth result: status=%d node=%s docker_id=%s container=%s tls_override=%t\n", status, node, docker_id, container, tls_override)
	return status, node, docker_id, container, tls_override
}

func get_id_from_uri(uri string, pattern string) string{
	var id string
	slice1 := strings.Split(uri, pattern)
	//log.Printf("get_id_from_uri: pattern=%s, slice1=%v\n", pattern, slice1)
	if len(slice1) > 1 {
		slice2 := strings.Split(slice1[1], "/")
		id = slice2[0]
		slice3 := strings.Split(id, "?")
		id = slice3[0]
	}else{
		id=""
	}
	//log.Printf("get_id_from_uri: id=%s\n", id)
	return id
}

func get_id_and_type(uri string) (id string, id_type string){
	patterns := []string {
		"/containers/json",
		"/containers/create",
		"/images",
		"/build",
		"/version",
		"/_ping",
	}

	found := false
	for i:=0; i < len(patterns); i++ {
		//if uri contains patterns[i]
		if strings.Contains(uri, patterns[i]) {
			found = true
		}
	}
	if found {
		id=""
		id_type = "None"
		log.Printf("id=%s, id_type=%s\n", id, id_type)
		return id, id_type
	}

	//1st: look for /containers/<id>/
	id = get_id_from_uri(uri, "/containers/")
	id_type="Container"
	if id == "" {
		//2nd: look for /exec/<id>/
		id = get_id_from_uri(uri, "/exec/")
		if id == "" {
			id_type = "None"
		}else {
			//id found in uri
			id_type = "Exec"
		}
	}
	log.Printf("id=%s, id_type=%s\n", id, id_type)
	return id, id_type
}


