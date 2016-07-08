package auth

import (
	"net/http"
	"strings"

	"conf"  		// my conf package
)


// returns auth=true/false, compute node name, container/exec id, container id,
// override tls flag is used in swarm case only
func DockerAuth(r *http.Request) (creds Creds) {

	// Make a call to get Shard and authenticate. This call is always sufficient for Swarm shard case.
	// For nova-docker case, AND exec calls, redis must be called to get container id then a second call to getHost is needed.

	//parse r.RequestURI for container id or exec id
	uri := r.RequestURI
	id, id_type := get_id_and_type(uri)
	//if type is Exec then get container id from redis cache to forward to getHost api
	var container_id string
	second_call_needed := false
	if id_type == "Exec" {
		//"Exec" AND nova-docker (to be determined based on 1st call response) ==> a 2nd call to getHost will be needed
		//1st call will use "NoneContainer", 2nd call will use the container_id to be retrieved from redis
		second_call_needed = true  // assume true for now until proven not to be nova-docker shard
		id_type = "None"
		//container_id = conf.RedisGet(id)
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
		Log.Printf("Auth result: status=%d. Trying again without the container_id\n", creds.Status)
		// try again without passing the container_id
		creds.Status, host = getHost(r, "NoneContainer")
 		if creds.Status != 200 {
			Log.Printf("Second Auth result: status=%d\n", creds.Status)
			return
		}
 	}

	if second_call_needed {
		id_type = "Exec"
		if ! host.Swarm {
			// nova-docker case, make 2nd call
			container_id = conf.RedisGet(id)
			creds.Status, host = getHost(r, container_id)
			if creds.Status != 200 {
				Log.Printf("Auth result: status=%d\n", creds.Status)
				return
			}
		}else{
			// swarm shard case, no need for 2nd call
			second_call_needed = false
		}
	}

	creds.Swarm_shard = host.Swarm
	creds.Node = host.Host + ":" + conf.GetDockerPort()
	creds.Container = host.Container_id
	creds.Reg_namespace = host.Namespace
	creds.Apikey = host.Apikey
	creds.Space_id =host.Space_id
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
		// this is required by swarm-auth
		//r.Header.Set("X-Auth-Token", GetNamespace(host.Space_id))
		r.Header.Set(conf.GetSwarmAuthHeader(), GetNamespace(host.Space_id))
		Log.Printf("Injected swarm-auth required header %s=%s", conf.GetSwarmAuthHeader(), GetNamespace(host.Space_id))
		if id_type == "Container" {
			//@@ This fix is needed especially for after ccsapi is fixed to not invoke swarm in getHost
			creds.Docker_id = id //creds.Container
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
	}

	// this check loop is important for prefixes that start with /containers or /exec and need to be excluded (return id_type="None")
	// other prefixes will lead to id_type="None" anyway
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


func GetNamespace(space_id string) (namespace string) {
	return "s" + space_id + "-default"
}
