package auth

import (
	"net/http"
	"strings"
	"strconv"

	"conf"  		// my conf package
)

// returns auth=true/false, compute node name, container/exec id, container id,
// override tls flag is used in swarm case only
func DockerAuth(r *http.Request) (creds Creds) {

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
		creds.Status, host = getHost(r, "NoneContainer")
	}else{
		creds.Status, host = getHost(r, container_id)
	}
	if creds.Status != 200 {
		Log.Printf("Auth result: status=%d\n", creds.Status)
		return
	}
	creds.Node = host.Host + ":" + conf.GetDockerPort()
	creds.Container = host.Container_id
	creds.Reg_namespace = host.Namespace
	creds.Apikey = host.Apikey
	creds.Space_id = host.Space_id
	creds.Orguuid = host.Orguuid
	creds.Userid = host.Userid
	//container id needs nova- prefix
	//exec id does not need a prefix
	if id_type == "Container" {
		creds.Docker_id = "nova-"+creds.Container   //append nova to the id returned from getHost
	}
	if id_type == "Exec" {
		creds.Docker_id = id  //use the exec id that came in the original req
	}
	if id_type == "None" {
		creds.Docker_id = ""
	}
	creds.Tls_override = false
	if host.Swarm {
		creds.Node = host.Mgr_host    //Mgr_host = host:port
		//insert space_id in the header to be forwarded
		r.Header.Set("X-Auth-Token", host.Space_id)
		if id_type == "Container" {
			creds.Docker_id = creds.Container
		}
		if id_type == "Exec" {
			creds.Docker_id = id  //use the exec id that came in the original req
		}
		if id_type == "None" {
			creds.Docker_id = ""
		}
		if !host.Swarm_tls{
			creds.Tls_override = true  // no tls for this outbound req regardless of proxy conf
		}
		if host.Host != ""{
			// go to swarm node directly if known
			creds.Node = host.Host + ":" + strconv.Itoa(conf.GetSwarmNodePort())
		}
	}
	Log.Printf("Auth result: status=%d node=%s docker_id=%s container=%s tls_override=%t reg_namespace=%s\n",
		creds.Status, creds.Node, creds.Docker_id, creds.Container, creds.Tls_override, creds.Reg_namespace)
	return
}

func get_id_from_uri(uri string, pattern string) string{
	var id string
	slice1 := strings.Split(uri, pattern)
	//Log.Printf("get_id_from_uri: pattern=%s, slice1=%v\n", pattern, slice1)
	if len(slice1) > 1 {
		slice2 := strings.Split(slice1[1], "/")
		id = slice2[0]
		slice3 := strings.Split(id, "?")
		id = slice3[0]
	}else{
		id=""
	}
	//Log.Printf("get_id_from_uri: id=%s\n", id)
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
		Log.Printf("id=%s, id_type=%s\n", id, id_type)
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
	Log.Printf("id=%s, id_type=%s\n", id, id_type)
	return id, id_type
}


