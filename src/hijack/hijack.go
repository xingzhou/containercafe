package main

import (
    "fmt"
    "net/http"
	"net/http/httputil"
	"log"
	"strconv"
	"io/ioutil"
	"os"
	"httphelper"  //my httphelper package
	"time"
	"strings"
	"net"
	"bytes"
	"sync"
	"auth"  // my auth package
	"conf"  // my conf package
	"crypto/tls"
	"encoding/json"
)

//TODO logging refinement
//TODO use req id generator using hashing fn

var glob_req_id = 0
var glob_req_id_mutex sync.Mutex

func get_req_id() int{
	glob_req_id_mutex.Lock()
	glob_req_id += 1  //this op should be in a critical section
	req_id := glob_req_id
	glob_req_id_mutex.Unlock()
	return req_id
}

func redirect_lowlevel(r *http.Request, body []byte, redirect_host string, redirect_resource_id string) (*http.Response, error, *httputil.ClientConn){
	//forward request to server
	//var c net.Conn
	//var err error
	//var buf *bufio.Reader = nil
	//var c_tls *tls.Conn
	var cc *httputil.ClientConn

	c , err := net.Dial("tcp", redirect_host)
	if err != nil {
		// handle error
		fmt.Printf("Error connecting to server %s, %v\n", redirect_host, err)
		return nil,err,nil
	}

	if conf.IsTlsOutbound() {
		c_tls := tls.Client(c, &tls.Config{InsecureSkipVerify : true})
		cc = httputil.NewClientConn(c_tls, nil)
	}else{
		cc = httputil.NewClientConn(c, nil)
	}

	req, _ := http.NewRequest(r.Method, "http://"+redirect_host+auth.RewriteURI(r.RequestURI, redirect_resource_id),
				bytes.NewReader(body))
	req.Header = r.Header
	//req.Host = redirect_host
	req.URL.Host = redirect_host

	resp, err := cc.Do(req)

	//defer resp.Body.Close()
	return resp, err, cc
}

func handler(w http.ResponseWriter, r *http.Request, redirect_host string, redirect_resource_id string) {
	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false
	resp_DOCKER := false
	//req_LOGS := false

	var err error = nil

	now := time.Now()
	req_id := get_req_id()
	fmt.Printf("------> req id: %d, time: %s\n", req_id, now)

	data, _ := httputil.DumpRequest(r, true)
	fmt.Printf("Request dump of %d bytes:\n", len(data))
	fmt.Printf("%s\n", string(data))

	body, _ := ioutil.ReadAll(r.Body)

	if (httphelper.IsUpgradeHeader(r.Header)) {
		fmt.Printf("@ Upgrade request detected\n")
		req_UPGRADE = true
	}

	//TODO ***** Filter framework for Interception of commands before forwarding request to server *****
	if is_container_attach_call(r.RequestURI) {
		//insert delay to allow for completion of container creation on the prior create command
		time.Sleep(15*time.Second)
	}
	//if is_container_logs_call(r.RequestURI) {
	//	fmt.Printf("@ Logs request detected\n")
	//	req_LOGS = true
	//}


	//resp, err := redirect(r, body, redirect_host)
	resp, err, cc := redirect_lowlevel(r, body, redirect_host, redirect_resource_id)
	if (err != nil) {
		fmt.Printf("Error in redirection... %v\n", err)
		//log.Fatal(err) //this would terminate the server
		fmt.Printf("___________________________ id: %d \n", req_id)
		return
	}

	//write out resp
	now = time.Now()
	fmt.Printf("<------ id: %d, time: %s\n", req_id, now)
	//data2, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("Response dump of %d bytes:\n", len(data2))
	//fmt.Printf("%s\n", string(data2))

	fmt.Printf("Resp Status: %s\n", resp.Status)
	httphelper.DumpHeader(resp.Header)

	//echo r.Header in resp_header
	//w.Header().Add("foo", "bar")   //ok
	httphelper.CopyHeader(w.Header(), resp.Header)

	//fmt.Printf(">> Debug forwarded response header BEGIN\n")
	//httphelper.DumpHeader(w.Header())
	//fmt.Printf("<< Debug forwarded response header END\n")

	if (httphelper.IsUpgradeHeader(resp.Header)) {
		fmt.Printf("@ Upgrade response detected\n")
		resp_UPGRADE = true
	}

	if httphelper.IsStreamHeader(resp.Header) {
		fmt.Printf("@ application/octet-stream detected\n")
		resp_STREAM = true
	}

	if httphelper.IsDockerHeader(resp.Header) {
		fmt.Printf("@ application/vnd.docker.raw-stream detected\n")
		resp_DOCKER = true
	}

	//TODO ***** Filter framework for Interception of commands before forwarding resp to client (1) *****
	//if req_LOGS {
	//	//insert streaming header in response to client
	//	w.Header().Set("Content-Type", "application/octet-stream")
	//}

	proto := strings.ToUpper(httphelper.GetHeader(resp.Header, "Upgrade"))
	if (req_UPGRADE || resp_UPGRADE) && (proto != "TCP") {
		fmt.Printf("@ Warning: will start hijack proxy loop although Upgrade proto %s is not TCP\n", proto)
	}

	if req_UPGRADE || resp_UPGRADE || resp_STREAM || resp_DOCKER {  // || req_LOGS{

		//resp header is sent first thing on hijacked conn
		w.WriteHeader(resp.StatusCode)

		fmt.Printf("@ starting tcp hijack proxy loop\n")
		httphelper.InitProxyHijack(w, cc, req_id, "TCP") // TCP is the only supported proto now
	}else {
		//If no hijacking, forward full response to client
		w.WriteHeader(resp.StatusCode)

		if resp.Body != nil {
			resp_body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Printf("Error: error in reading server response body\n")
			}else {
				//TODO ***** Filter framework for Interception of commands before returning result to client (2) *****
				//Check if Redis caching is required
				//if request uri contains "/container/" and "/exec" then store in Redis the returned exec id (in resp body) and container id (in uri)
				if is_container_exec_call(r.RequestURI){
					container_id := strip_nova_prefix(redirect_resource_id)
					exec_id := get_exec_id_from_response(resp_body)
					if exec_id == ""{
						fmt.Printf("Error: error in retrieving exec id from response body\n")
					}else {
						conf.RedisSetExpire(exec_id, container_id, 60*60)
					}
				}

				//Printout the response body
				if strings.ToLower(httphelper.GetHeader(resp.Header, "Content-Type")) == "application/json" {
					httphelper.PrintJson(resp_body)
				}else{
					//fmt.Printf("Received %d bytes\n", len(resp_body))
					fmt.Printf("\n%s\n", string(resp_body))
				}

				//forward server response to calling client
				fmt.Fprintf(w, "%s", resp_body)
			}
		}else {
			fmt.Printf("\n")
			fmt.Fprintf(w, "\n")
		}
	}

	fmt.Printf("___________________________ id : %d \n", req_id)
	return
}

//http proxy forwarding with hijack support handler
func endpoint_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("@ endpoint_handler triggered, URI: %s\n", r.RequestURI)
	ok, node, docker_id := auth.Auth(r)  // ok=true/false, node=host:port, docker_id=url resource id understood by docker
	if !ok {
		return
	}
	handler(w, r, node, docker_id)
}

//Return 404 for all non-supported URIs
func no_endpoint_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("@ no_endpoint_handler triggered, URI: %s, returning error 404\n", r.RequestURI)
	//w.WriteHeader(404)
	http.NotFound(w, r)
}

func health_endpoint_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("@ health_endpoint_handler triggered, URI: %s\n", r.RequestURI)
	fmt.Fprintf(w,"hjproxy up\n")
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

	fmt.Printf("@ get_exec_id_from_response: json=%s\n", body)
	err := json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("@ get_exec_id_from_response: error=%v", err)
		return ""
	}
	fmt.Printf("@ get_exec_id_from_response: Id=%s\n", resp.Id)
	return resp.Id
}

func main() {
	fmt.Printf("Hijack proxy microservice started...\n")
	conf.LoadEnv()
	listen_port := conf.GetDefaultListenPort()
	//parse args
	nargs := len(os.Args)
	if nargs > 1 {
		listen_port, _=strconv.Atoi(os.Args[1])
	}
	fmt.Printf("@ Listening on port %d\n", listen_port)
	if nargs > 2 {
		conf.Default_redirect_host = os.Args[2]
	}
	fmt.Printf("@ Default upstream server: %s\n", conf.Default_redirect_host)
	if nargs > 3 {
		if strings.ToLower(os.Args[3]) == "true" {
			conf.SetTlsInbound(true)
			conf.SetTlsOutbound(true)
		}
		if strings.ToLower(os.Args[3]) == "false" {
			conf.SetTlsInbound(false)
			conf.SetTlsOutbound(true)
		}
	}
	fmt.Printf("@ tls setup: Inbound=%t, Outbound=%t\n", conf.IsTlsInbound(), conf.IsTlsOutbound())

	//register handlers for supported url paths, can't register same path twice

	//Healthcheck handler
	http.HandleFunc("/hjproxy/", health_endpoint_handler)

	//Rely on NGINX to route accepted url paths
	http.HandleFunc("/", endpoint_handler)
	/*
	http.HandleFunc("/v1.18/containers/", endpoint_handler)			//  /<v>/containers/<id>/logs ... hijack
	http.HandleFunc("/v1.17/containers/", endpoint_handler)			//  /<v>/containers/<id>/attach ... hijack
	http.HandleFunc("/v1.20/containers/", endpoint_handler)
	http.HandleFunc("/v3/containers/", endpoint_handler)
	http.HandleFunc("/v1.18/exec/", endpoint_handler)				//  /<v>/exec/<id>/start ... hijack
	http.HandleFunc("/v1.17/exec/", endpoint_handler)
	http.HandleFunc("/v1.20/exec/", endpoint_handler)
	http.HandleFunc("/v3/exec/", endpoint_handler)
	*/

	//Disable trap handler for non-supported url paths. Rely on NGINX to route accepted url paths
	//http.HandleFunc("/", no_endpoint_handler)

	//init server on any interface + listen_port

	var err error
	if conf.IsTlsInbound() {
		err = http.ListenAndServeTLS(":"+strconv.Itoa(listen_port), conf.GetCertFile(), conf.GetKeyFile(), nil)
	} else {
		err = http.ListenAndServe(":"+strconv.Itoa(listen_port), nil)
	}

	//print something and exit on fatal error
	if err != nil {
		log.Fatal("Hijack microservice aborting because could not start ListenAndServe: ", err)
	}
}

