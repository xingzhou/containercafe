//docker and swarm handler
//
package handler

import (
    "fmt"
    "net/http"
	"net/http/httputil"
	"io/ioutil"
	"time"
	"strings"
	"encoding/json"
	"net/url"

	"limit"  //my limits package
	"httphelper"  //my httphelper package
	"auth"  // my auth package
	"conf"  // my conf package
)

// supported docker api uri patterns
var dockerPatterns = []string {
	"/containers/",
	"/images/",
	"/exec/",
	"/version",
	"/auth",
	"/_ping",
	"/build", 	// ?? buildsrvc??
	"/commit", 	// ?? create img from  container content
	"/info", 	// ??  system wide info
	"/events",	// ??
}

// http proxy forwarding with hijack support
// handler for docker/swarm
func DockerEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	Log.Printf("------> DockerEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// check if uri pattern is accepted
	if ! IsSupportedPattern(r.RequestURI, dockerPatterns){
		Log.Printf("Docker pattern not accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
		NoEndpointHandler(w, r)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	// Call Auth interceptor
	// ok=true/false, node=host:port,
	// docker_id=resource id from url mapped to id understood by docker
	// container != docker_id in exec case
	// tls_override is true when swarm master does not support tls
	status, node, docker_id, container, tls_override, reg_namespace := auth.DockerAuth(r)
	if status != 200 {
		Log.Printf("Authentication failed for req_id=%s status=%d", req_id, status)
		if status == 401 {
			NotAuthorizedHandler(w,r)
		}else{
			ErrorHandler(w,r,status)
		}
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
    Log.Printf("Authentication succeeded for req_id=%s status=%d", req_id, status)


	data, _ := httputil.DumpRequest(r, true)
	Log.Printf("Request dump req_id=%s req_length=%d:\n%s", req_id, len(data), string(data))
	body, _ := ioutil.ReadAll(r.Body)

	// Check for attempt to access not-allowed image
	img := ""
	check_img := false
	// First check - run call
	if is_container_create_call(r.RequestURI){
		// extract image
		img = get_image_from_container_create(body)
		check_img = true

		// inject X-Registry-Auth header
		InjectRegAuthHeader(r)
	}
	// Second check - pull call
	if is_image_create_call(r.RequestURI) {
		// extract image
		img = get_image_from_image_create(r.RequestURI)
		check_img = true
	}
	// Third check - push call
	if is_image_push_call(r.RequestURI) {
		// extract image
		img = get_image_from_image_push(r.RequestURI)
		check_img = true
	}
	if is_image_inspect_call(r.RequestURI){
		// extract image
		img = get_image_from_image_inspect(r.RequestURI)
		check_img = true
	}
	if check_img && !is_img_valid(img, reg_namespace){
		// check that image name is valid for this user
		Log.Printf("Not allowed to access image img=%s namespace=%s req_id=%s", img, reg_namespace, req_id)
		//NotAuthorizedHandler(w, r)
		ErrorHandlerWithMsg(w, r, 500, "Not allowed to access image")
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	// intercept image calls and direct to registry
	if is_image_list_call(r.RequestURI){
		invoke_reg_list(w, r, reg_namespace, req_id)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
	if is_image_inspect_call(r.RequestURI){
		invoke_reg_inspect(w, r, img, reg_namespace, req_id)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
	if is_image_rmi_call(r.RequestURI, r.Method){
		invoke_reg_rmi(w, r, reg_namespace, req_id)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Call conn limiting interceptor(s) pre-processing
	if !limit.OpenConn(container, conf.GetMaxContainerConn()) {
		Log.Printf("Max conn limit reached for container...aborting request")
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
	if !limit.OpenConn(node, conf.GetMaxNodeConn()) {
		Log.Printf("Max conn limit reached for host node...aborting request")
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Handle request
	dockerHandler(w, r, body, node, docker_id, req_id, tls_override)

	//Call conn limiting interceptor(s) post-processing, to decrement conn count(s)
	limit.CloseConn(container, conf.GetMaxContainerConn())
	limit.CloseConn(node, conf.GetMaxNodeConn())

	Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
}

// private handler processing
func dockerHandler(w http.ResponseWriter, r *http.Request, body []byte, redirect_host string,
	redirect_resource_id string, req_id string, tls_override bool) {

	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false
	resp_DOCKER := false
	req_LOGS := false

	var err error = nil

	//***** Filter req/headers here before forwarding request to server *****

	if (httphelper.IsUpgradeHeader(r.Header)) {
		Log.Printf("@ Upgrade request detected\n")
		req_UPGRADE = true
	}
	if is_container_logs_call(r.RequestURI) {
		Log.Printf("@ Logs request detected\n")
		req_LOGS = true
	}

	maxRetries := 1
	backOffTimeout := 0
	if is_container_attach_call(r.RequestURI) {
		//insert delay to allow for completion of container creation on the prior create command
		//time.Sleep(15*time.Second)
		maxRetries = conf.GetMaxRetries()
		backOffTimeout = conf.GetBackOffTimeout()
	}

	var (resp *http.Response
		cc *httputil.ClientConn
	)
	for i:=0; i<maxRetries; i++ {
		resp, err, cc = redirect(r, body, redirect_host, redirect_resource_id,
			dockerRewriteUri, tls_override)
		if err == nil {
			break
		}
		Log.Printf("redirect failed retry=%d req_id=%s", i, req_id)
		if (i+1) < maxRetries {
			Log.Printf("will sleep before retry secs=%d req_id=%s", backOffTimeout, req_id)
			time.Sleep( time.Duration(backOffTimeout) * time.Second)
		}
	}
	if (err != nil) {
		Log.Printf("Error in redirection, will abort req_id=%s err=%v\n", req_id, err)
		return
	}

	//write out resp
	//now = time.Now()
	Log.Printf("<------ req_id=%s\n", req_id)
	//data2, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("Response dump of %d bytes:\n", len(data2))
	//fmt.Printf("%s\n", string(data2))

	Log.Printf("Resp Status: %s\n", resp.Status)
	Log.Print( httphelper.DumpHeader(resp.Header) )

	httphelper.CopyHeader(w.Header(), resp.Header)

	if (httphelper.IsUpgradeHeader(resp.Header)) {
		Log.Printf("@ Upgrade response detected\n")
		resp_UPGRADE = true
	}
	if httphelper.IsStreamHeader(resp.Header) {
		Log.Printf("@ application/octet-stream detected\n")
		resp_STREAM = true
	}
	if httphelper.IsDockerHeader(resp.Header) {
		Log.Printf("@ application/vnd.docker.raw-stream detected\n")
		resp_DOCKER = true
	}

	//***** Filter framework for Interception of commands before forwarding resp to client (1) *****

	proto := strings.ToUpper(httphelper.GetHeader(resp.Header, "Upgrade"))
	if (req_UPGRADE || resp_UPGRADE) && (proto != "TCP") {
		Log.Printf("Warning: will start hijack proxy loop although Upgrade proto %s is not TCP\n", proto)
	}

	if req_UPGRADE || resp_UPGRADE || resp_STREAM || resp_DOCKER || req_LOGS{
		//resp header is sent first thing on hijacked conn
		w.WriteHeader(resp.StatusCode)

		Log.Printf("starting tcp hijack proxy loop req_id=%s", req_id)
		httphelper.InitProxyHijack(w, cc, req_id, "TCP") // TCP is the only supported proto now
		return
	}
	//If no hijacking, forward full response to client
	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		Log.Printf("\n")
		fmt.Fprintf(w, "\n")
		return
	}

	_DOCKER_CHUNKED_READ_ := false   // new feature flag

	if _DOCKER_CHUNKED_READ_ {
		//new code to test
		//defer resp.Body.Close()   // causes this method to not return to caller IF closing while there is still data in Body!
		chunkedRWLoop(resp, w, req_id)

		// TODO extract exec id from resp

	}else {
		resp_body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			Log.Printf("Error: error in reading server response body")
			fmt.Fprint(w, "error in reading server response body\n")
			return
		}

		//***** Filter framework for Interception of commands before returning result to client (2) *****
		//Check if Redis caching is required
		//if request uri contains "/container/" and "/exec" then store in Redis the returned exec id (in resp body) and container id (in uri)
		if is_container_exec_call(r.RequestURI) {
			container_id := strip_nova_prefix(redirect_resource_id)
			exec_id := get_exec_id_from_response(resp_body)
			if exec_id == "" {
				Log.Printf("Error: error in retrieving exec id from response body")
			}else {
				conf.RedisSetExpire(exec_id, container_id, 60*60)
			}
		}

		//Printout the response body
		bodystr := "Dump Body:\n"
		if strings.ToLower(httphelper.GetHeader(resp.Header, "Content-Type")) == "application/json" {
			bodystr += httphelper.PrettyJson(resp_body)
		}else {
			bodystr += string(resp_body)+"\n"
		}
		Log.Println(bodystr)

		//forward server response to calling client
		fmt.Fprintf(w, "%s", resp_body)
	}
	return
}

//Convert /v*/containers/id/*  to  /<docker_api_ver>/containers/<redirect_resource_id>/*
//Convert /v*/exec/id/*  to  /<docker_api_ver>/exec/<redirect_resource_id>/*
func dockerRewriteUri(reqUri string, redirect_resource_id string)(redirectUri string){
	sl := strings.Split(reqUri, "/")
	if redirect_resource_id == "" {
		//supports /v../containers/json  /v../build  /v../build?foo=bar
		//redirectURI = conf.GetDockerApiVer()+"/"+sl[2]+"/"+sl[3]
		redirectUri = conf.GetDockerApiVer()
		for i:=2; i < len(sl); i++ {
			redirectUri += "/" + sl[i]
		}
	}else {
		//redirectURI = conf.GetDockerApiVer()+"/"+sl[2]+"/"+redirect_resource_id+"/"+sl[4]
		redirectUri = conf.GetDockerApiVer()+"/"+sl[2]+"/"+redirect_resource_id
		for i:=4; i < len(sl); i++ {
			redirectUri += "/" + sl[i]
		}
		//what if there is ?foo=bar in last slice and last slice is resource_id e.g., DELETE /v/containers/123?foo=bar
		if len(sl) <= 4 {
			sl2 := strings.Split(sl[len(sl)-1], "?")
			if len(sl2) > 1{
				redirectUri += "?" + sl2[1]
			}
		}
	}
	Log.Printf("dockerRewriteURI: '%s' --> '%s'\n", reqUri, redirectUri)
	return redirectUri
}



func strip_nova_prefix(id string) string{
	return strings.TrimPrefix(id, "nova-")
}

func get_exec_id_from_response(body []byte) string{
	type Resp struct {
		Id  		string
		Warnings 	[]string
	}
	var resp Resp

	Log.Printf("get_exec_id_from_response: json=%s\n", body)
	err := json.Unmarshal(body, &resp)
	if err != nil {
		Log.Printf("get_exec_id_from_response: error=%v", err)
		return ""
	}
	Log.Printf("get_exec_id_from_response: Id=%s\n", resp.Id)
	return resp.Id
}

//////////////////////////// image names extraction and validation ops

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
	// Ex: POST /images/registry.acme.com:5000/test/push HTTP/1.1
	sl := strings.Split(reqUri, "/")
	if len(sl) < 4 {
		// err
		img=""
		return
	}
	img = sl[2]
	for i:=3; i< len(sl)-1; i++ {
		img += "/" + sl[i]
	}
	return
}

func get_image_from_image_inspect(reqUri string) (img string){
	// Ex: POST /images/registry.acme.com:5000/test/json HTTP/1.1
	sl := strings.Split(reqUri, "/")
	if len(sl) < 4 {
		// err
		img=""
		return
	}
	img = sl[2]
	for i:=3; i< len(sl)-1; i++ {
		img += "/" + sl[i]
	}
	return
}

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

//////////////////////////// Check request URI for a certain call pattern

//return true if it is /<v>/containers/<id>/exec api call
func is_container_exec_call(uri string) bool {
	if strings.Contains(uri, "/containers/") && strings.Contains(uri, "/exec") {
		return true
	}else{
		return false
	}
}

func is_container_attach_call(uri string) bool {
	if strings.Contains(uri, "/containers/") && strings.Contains(uri, "/attach") {
		return true
	}else{
		return false
	}
}

func is_container_logs_call(uri string) bool {
	if strings.Contains(uri, "/containers/") && strings.Contains(uri, "/logs") {
		return true
	}else{
		return false
	}
}

func is_container_create_call(uri string) bool {
	if strings.Contains(uri, "/containers/create") {
		return true
	}else{
		return false
	}
}

//image pull call
func is_image_create_call(uri string) bool {
	if strings.Contains(uri, "/images/create") {
		return true
	}else{
		return false
	}
}

func is_image_push_call(uri string) bool {
	if strings.Contains(uri, "/images/") && strings.Contains(uri, "/push"){
		return true
	}else{
		return false
	}
}

func is_image_list_call(uri string) bool{
	if strings.Contains(uri, "/images/json"){
		return true
	}else{
		return false
	}
}

func is_image_inspect_call(uri string) bool{
	if strings.Contains(uri, "/images/") && strings.Contains(uri, "/json") && !strings.Contains(uri, "/images/json"){
		return true
	}else{
		return false
	}
}

func is_image_rmi_call(uri, method string) bool{
	if strings.Contains(uri, "/images/") && method == "DELETE" {
		return true
	}else{
		return false
	}
}
