package handler

import (
	"net/http"
	"net/url"
	"io"
	"io/ioutil"
	"strings"
	"fmt"
	"encoding/base64"
	"encoding/json"

	"auth"
	"conf"
	"httphelper"
)

func InjectRegAuthHeader(r *http.Request, creds auth.Creds) {
	//Use admin base64 encoded psswd
	//tok := conf.GetRegAuthToken()

	//create X-Registry-Auth object out of apikey in creds
	// {"username":"admin","password":"230189","auth":"","email":"swarm@dev.test","serveraddress":"registry-ice-dev-test.stage1.ng.bluemix.net"}
	auth_str := fmt.Sprintf("{\"username\":\"apikey\",\"password\":\"%s\",\"auth\":\"\",\"email\":\"swarm@dev.test\",\"serveraddress\":\"%s\"}", creds.Apikey, conf.GetRegLocation())
    auth_bytes := []byte(auth_str)
	tok := base64.StdEncoding.EncodeToString( auth_bytes )

	Log.Printf("InjectRegAuthHeader:  auth_str=%s  tok=%s", auth_str, tok)
	r.Header.Set("X-Registry-Auth", tok)
}

func GetRegistryApiHost() (host string){
	// get service host from Consul
	service := "registry-api"
	host = conf.GetServiceHost(service)
	if (host == ""){
		Log.Printf("Failed to get Registry API host  service=%s", service)
	}
	return
}

func AddCredsHeaders(req *http.Request, creds auth.Creds){
	req.Header.Add("namespace", creds.Reg_namespace)
	req.Header.Add("apikey", creds.Apikey)
	req.Header.Add("orguuid", creds.Orguuid)
	req.Header.Add("spaceuuid", creds.Space_id)
	req.Header.Add("userid", creds.Userid)

	req.Header.Add("Accept", "application/json")
}

// call internal registry api server to get image metadata
func invoke_reg_inspect(w http.ResponseWriter, r *http.Request, img string, creds auth.Creds, req_id string){
	host := GetRegistryApiHost()
	if (host == ""){
		ErrorHandlerWithMsg(w, r, 500, "Failed to get Registry API host")
		return
	}

	//Call service endpoint
	url := "http://"+host+"/v1/imageJson?imageName="+img
	Log.Printf("Will call Registry... url=%s req_id=%s", url, req_id)
	client := &http.Client{}
	req, _ := http.NewRequest("GET",url, nil)
	AddCredsHeaders(req, creds)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Registry status code: %d", resp.StatusCode)
	Log.Printf("Registry response: %s", httphelper.PrettyJson(body) )

	//send response back to client
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, string(body))

	return
}

//implement image list by invoking search api of Containers registry
//return json to docker cli
//DockerHandler will print req processing completion message and exit right after this method
func invoke_reg_list(w http.ResponseWriter, r *http.Request, creds auth.Creds, req_id string){
	//Not recommended approach. These urls go to registry.ng.bluemix.net
	//ns_url := "http://" + conf.GetRegLocation() + "/v1/namespaces/" + namespace
	//lib_url := "http://" + conf.GetRegLocation() + "/v1/namespaces/library"

	// Recommended approach using internal api exposed by registry microservice
	// get service host from Consul
	host := GetRegistryApiHost()
	if (host == ""){
		ErrorHandlerWithMsg(w, r, 500, "Failed to get Registry API host")
		return
	}

	//Call service endpoint
	url := "http://" + host + "/v1/imageList/" + creds.Reg_namespace
	Log.Printf("Will call Registry... url=%s req_id=%s", url, req_id)
	client := &http.Client{}
	req, _ := http.NewRequest("GET",url, nil)
	AddCredsHeaders(req, creds)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Registry status code: %d", resp.StatusCode)
	Log.Printf("Registry response: %s", httphelper.PrettyJson(body))

	//send response back to client
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, string(body))
	return
}

func invoke_reg_rmi(w http.ResponseWriter, r *http.Request, img string, creds auth.Creds, req_id string){
	//img valid format:  host/namespace/name:tag
	sl := strings.Split(img, "/")
	if len(sl) != 3  {
		Log.Printf("Not valid image name  img=%s req_id=%s", img, req_id)
		ErrorHandlerWithMsg(w, r, 500, "Not valid image name")
		return
	}
	reg := sl[0]
    namespace := sl[1]
	img_name_and_tag := strings.Split(sl[2], ":")
	if len(img_name_and_tag) != 2 {
		Log.Printf("Not valid image name  img=%s req_id=%s", img, req_id)
		ErrorHandlerWithMsg(w, r, 500, "Not valid image name")
		return
	}
	img_name := img_name_and_tag[0]
	img_tag := img_name_and_tag[1]

	// Construct the request to registry.
    url := fmt.Sprintf("https://%s/v1/repositories/%s/%s/tags/%s", reg, namespace, img_name, img_tag)
	Log.Printf("Will call Registry... url=%s req_id=%s", url, req_id)
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE",url, nil)
	//TODO: when rmi is implemented by reg microservice 	AddCredsHeaders(req, creds)
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("apikey", creds.Apikey)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Registry status code: %d", resp.StatusCode)
	Log.Printf("Registry response: %s", httphelper.PrettyJson(body))

	//send response back to client
	body_str := ""
	if resp.StatusCode == 200{
		body_str = fmt.Sprintf("[{\"Deleted\":\"%s\"}]", img)
	}else{
		body_str = string(body)
	}
	w.WriteHeader(resp.StatusCode)
	io.WriteString(w, body_str)
	return
}

//////////////////////////// image names extraction and validation ops

func is_img_valid(img string, namespace string) bool{
	if img == "" {
		return true
	}

	// img general format is reg_host/namesapce/imgname:tag or reg_host/imgname:tag
	sl := strings.Split(img, "/")
	if len(sl) <= 2 {
		// public/lib image
		return true
	}

	// we have an img with a namespace
	if !strings.Contains(sl[0], ".bluemix.net") {
		//Dckerhub or other reg image --> OK
		return true
	}

	// we have a Containers reg img with namespace
	// limit access only to this environment's registry
	if namespace == sl[1] && conf.GetRegLocation() == sl[0]{
		return true
	}
	return false
}

func get_image_from_container_create(body []byte) (img string){
	// look for "Image":"..."
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil{
		Log.Printf("get_image_from_container_create: error in json unmarshalling, err=%v", err)
		return
	}
	m := f.(map[string]interface{})
	for k, v := range m {
		if (k == "Image") {
			img = v.(string)
			Log.Printf("get_image_from_container_create: found img=%s", img)
			return
		}
	}
	Log.Print("get_image_from_container_create: did not find Image in json body")
	return
}

func getImageFullnameFromVars(vars map[string]string)(fullname string){
	if _, ok := vars["{img}"]; !ok {
		return
	}
	fullname = vars["{img}"]
	if ns, ok := vars["{ns}"]; ok {
		fullname = ns + "/" + fullname
	}
	if reg, ok := vars["{reg}"]; ok {
		fullname = reg + "/" + fullname
	}
	return
}

func get_image_from_image_create(reqUri string) (img string){
	//look for ?fromImage=...&registry=...
	u, err := url.Parse(reqUri)
	if err != nil {
		Log.Print(err)
		return
	}
	q := u.Query()  // q is map[string][]string
	fromImage := q.Get("fromImage")
	registry := q.Get("registry)")
	if (registry == ""){
		img = fromImage
		return
	}
	if strings.Contains(fromImage, registry){
		img = fromImage
		return
	}
	img = registry+"/"+fromImage
	return
}

func get_image_from_image_push(reqUri string) (img string){
	// Ex: POST /v1.20/images/registry.acme.com:5000/test/push HTTP/1.1
	sl := strings.Split(reqUri, "/")
	if len(sl) < 5 {
		// err
		img=""
		return
	}
	img = sl[3]
	for i:=4; i< len(sl)-1; i++ {
		img += "/" + sl[i]
	}
	return
}

func get_image_from_image_inspect(reqUri string) (img string){
	// Ex: POST /v1.20/images/registry.acme.com:5000/test/json HTTP/1.1
	sl := strings.Split(reqUri, "/")
	if len(sl) < 5 {
		// err
		img=""
		return
	}
	img = sl[3]
	for i:=4; i< len(sl)-1; i++ {
		img += "/" + sl[i]
	}
	return
}

func get_image_from_image_rmi(reqUri string) (img string){
	// Ex: DELETE /v1.20/images/registry.acme.com:5000/namespace/test:latest HTTP/1.1
	sl := strings.Split(reqUri, "/")
	if len(sl) < 6 {
		// err return ""
		return
	}
	img = sl[3]
	for i:=4; i< len(sl); i++ {
		img += "/" + sl[i]
	}
	return
}
