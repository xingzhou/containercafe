package handler

import (
	"net/http"
	"io"
	"io/ioutil"

	"conf"
)

func InjectRegAuthHeader(r *http.Request) {
	tok := conf.GetRegAuthToken()
	r.Header.Set("X-Registry-Auth", tok)
}

// call internal registry api server to get image metadata
func invoke_reg_inspect(w http.ResponseWriter, r *http.Request, img string, namespace string, req_id string){
	// get service host from Consul
	service := "registry-api-external"
	host := conf.GetServiceHost(service)
	if (host == ""){
		Log.Printf("Failed to get host  service=%s  req_id=%s", service, req_id)
		ErrorHandlerWithMsg(w, r, 500, "Failed to get Registry API host")
		return
	}

	//Call service endpoint
	url := "http://"+host+"/v1/imageJson?imageName="+img
	Log.Printf("Will call Registry... url=%s req_id=%s", url, req_id)
	client := &http.Client{}
	req, _ := http.NewRequest("GET",url, nil)
	req.Header.Add("Accept", "application/json")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Registry status code: %d", resp.StatusCode)
	Log.Printf("Registry response: %s", string(body))

	//send response back to client
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, string(body))

	return
}

//implement image list by invoking search api of Containers registry
//return json to docker cli
//DockerHandler will print req processing completion message and exit right after this method
func invoke_reg_list(w http.ResponseWriter, r *http.Request, namespace string, req_id string){
	//TODO
	//1- call internal registry api server to get namespace images
	//ns_url := "http://" + conf.GetRegLocation() + "/v1/namespaces/" + namespace

	//2- call to get library images
	//lib_url := "http://" + conf.GetRegLocation() + "/v1/namespaces/library"

	//3- TODO get images metadata (id, size, creation date, tags?)
	ErrorHandlerWithMsg(w, r, 404, "Registry list API not supported yet... Stay tuned")
	return
}

func invoke_reg_rmi(w http.ResponseWriter, r *http.Request, namespace string, req_id string){
	//TODO
}
