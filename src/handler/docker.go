//docker and swarm handler
//
package handler

import (
    "fmt"
    "net/http"
	"net/http/httputil"
	"log"
	"io/ioutil"
	"time"
	"strings"
	"encoding/json"

	"limit"  //my limits package
	"httphelper"  //my httphelper package
	"auth"  // my auth package
	"conf"  // my conf package
)

// http proxy forwarding with hijack support
// handler for docker/swarm
func DockerEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	log.Printf("------> DockerEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// Call Auth interceptor
	// ok=true/false, node=host:port,
	// docker_id=resource id from url mapped to id understood by docker
	// container != docker_id in exec case
	// tls_override is true when swarm master does not support tls
	ok, node, docker_id, container, tls_override := auth.DockerAuth(r)
	if !ok {
		log.Printf("Authentication failed for req_id=%s", req_id)
		log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Call conn limiting interceptor(s) pre-processing
	if !limit.OpenConn(container, conf.GetMaxContainerConn()) {
		log.Printf("Max conn limit reached for container...aborting request")
		log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
	if !limit.OpenConn(node, conf.GetMaxNodeConn()) {
		log.Printf("Max conn limit reached for host node...aborting request")
		log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Handle request
	dockerHandler(w, r, node, docker_id, req_id, tls_override)

	//Call conn limiting interceptor(s) post-processing, to decrement conn count(s)
	limit.CloseConn(container, conf.GetMaxContainerConn())
	limit.CloseConn(node, conf.GetMaxNodeConn())

	log.Printf("------ Completed processing of request req_id=%s\n", req_id)
}

// private handler processing
func dockerHandler(w http.ResponseWriter, r *http.Request, redirect_host string,
	redirect_resource_id string, req_id string, tls_override bool) {

	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false
	resp_DOCKER := false
	req_LOGS := false

	var err error = nil

	data, _ := httputil.DumpRequest(r, true)
	log.Printf("Request dump of %d bytes:\n%s", len(data), string(data))

	body, _ := ioutil.ReadAll(r.Body)

	//***** Filter req/headers here before forwarding request to server *****

	if (httphelper.IsUpgradeHeader(r.Header)) {
		log.Printf("@ Upgrade request detected\n")
		req_UPGRADE = true
	}
	if is_container_logs_call(r.RequestURI) {
		log.Printf("@ Logs request detected\n")
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
		log.Printf("redirect retry=%d failed", i)
		if (i+1) < maxRetries {
			log.Printf("will sleep secs=%d before retry", backOffTimeout)
			time.Sleep( time.Duration(backOffTimeout) * time.Second)
		}
	}
	if (err != nil) {
		log.Printf("Error in redirection, will abort req_id=%s ... err=%v\n", req_id, err)
		return
	}

	//write out resp
	//now = time.Now()
	log.Printf("<------ req_id=%s\n", req_id)
	//data2, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("Response dump of %d bytes:\n", len(data2))
	//fmt.Printf("%s\n", string(data2))

	log.Printf("Resp Status: %s\n", resp.Status)
	log.Print( httphelper.DumpHeader(resp.Header) )

	httphelper.CopyHeader(w.Header(), resp.Header)

	if (httphelper.IsUpgradeHeader(resp.Header)) {
		log.Printf("@ Upgrade response detected\n")
		resp_UPGRADE = true
	}
	if httphelper.IsStreamHeader(resp.Header) {
		log.Printf("@ application/octet-stream detected\n")
		resp_STREAM = true
	}
	if httphelper.IsDockerHeader(resp.Header) {
		log.Printf("@ application/vnd.docker.raw-stream detected\n")
		resp_DOCKER = true
	}

	//TODO ***** Filter framework for Interception of commands before forwarding resp to client (1) *****

	proto := strings.ToUpper(httphelper.GetHeader(resp.Header, "Upgrade"))
	if (req_UPGRADE || resp_UPGRADE) && (proto != "TCP") {
		log.Printf("Warning: will start hijack proxy loop although Upgrade proto %s is not TCP\n", proto)
	}

	if req_UPGRADE || resp_UPGRADE || resp_STREAM || resp_DOCKER || req_LOGS{
		//resp header is sent first thing on hijacked conn
		w.WriteHeader(resp.StatusCode)

		log.Printf("starting tcp hijack proxy loop\n")
		httphelper.InitProxyHijack(w, cc, req_id, "TCP") // TCP is the only supported proto now
		return
	}
	//If no hijacking, forward full response to client
	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		log.Printf("\n")
		fmt.Fprintf(w, "\n")
		return
	}
	//TODO chunked reads
	resp_body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Printf("Error: error in reading server response body\n")
		return
	}

	//TODO ***** Filter framework for Interception of commands before returning result to client (2) *****
	//Check if Redis caching is required
	//if request uri contains "/container/" and "/exec" then store in Redis the returned exec id (in resp body) and container id (in uri)
	if is_container_exec_call(r.RequestURI){
		container_id := strip_nova_prefix(redirect_resource_id)
		exec_id := get_exec_id_from_response(resp_body)
		if exec_id == ""{
			log.Printf("Error: error in retrieving exec id from response body\n")
		}else {
			conf.RedisSetExpire(exec_id, container_id, 60*60)
		}
	}

	//Printout the response body
	if strings.ToLower(httphelper.GetHeader(resp.Header, "Content-Type")) == "application/json" {
		httphelper.PrintJson(resp_body)
	}else{
		log.Printf("\n%s\n", string(resp_body))
	}

	//forward server response to calling client
	fmt.Fprintf(w, "%s", resp_body)
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
	log.Printf("dockerRewriteURI: '%s' --> '%s'\n", reqUri, redirectUri)
	return redirectUri
}

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

func strip_nova_prefix(id string) string{
	return strings.TrimPrefix(id, "nova-")
}

func get_exec_id_from_response(body []byte) string{
	type Resp struct {
		Id  		string
		Warnings 	[]string
	}
	var resp Resp

	log.Printf("get_exec_id_from_response: json=%s\n", body)
	err := json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("get_exec_id_from_response: error=%v", err)
		return ""
	}
	log.Printf("get_exec_id_from_response: Id=%s\n", resp.Id)
	return resp.Id
}

/*
func redirect_lowlevel(r *http.Request, body []byte, redirect_host string, redirect_resource_id string) (*http.Response, error, *httputil.ClientConn){
	//forward request to server
	var cc *httputil.ClientConn

	c , err := net.Dial("tcp", redirect_host)
	if err != nil {
		// handle error
		log.Printf("Error connecting to server %s, %v\n", redirect_host, err)
		return nil,err,nil
	}

	if conf.IsTlsOutbound() && !conf.GetTlsOutboundOverride(){
		cert, er := tls.LoadX509KeyPair(conf.GetClientCertFile(), conf.GetClientKeyFile())
		if er != nil {
			log.Printf("Error loading client key pair, %v\n", er)
			return nil,err,nil
		}
		c_tls := tls.Client(c, &tls.Config{InsecureSkipVerify : true, Certificates : []tls.Certificate{cert}})
		cc = httputil.NewClientConn(c_tls, nil)
	}else{
		cc = httputil.NewClientConn(c, nil)

		//The override is for the current request being processed only
		//The override is a directive received from ccsapi getHost, for a swarm request when swarm master does not support tls
		conf.SetTlsOutboundOverride(false)
	}

	req, _ := http.NewRequest(r.Method, "http://"+redirect_host+auth.RewriteURI(r.RequestURI, redirect_resource_id),
				bytes.NewReader(body))
	req.Header = r.Header
	//req.Host = redirect_host
	req.URL.Host = redirect_host

	log.Println("will forward request to server...")
	resp, err := cc.Do(req)

	//defer resp.Body.Close()
	return resp, err, cc
}
*/
