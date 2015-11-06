//kubernetes handler
//
package handler

import (
	"net/http"
	"fmt"
	"net/http/httputil"
	"io/ioutil"
	"time"
	"strings"

	"httphelper"
	"conf"
	"auth"
)

// supported Kubernetes api uri prefix patterns
var kubePatterns = []string {
	"/api/v1/namespaces/",
	"/api/v1/watch/namespaces/",
	"/api/v1/proxy/namespaces/",
	"/api",
	"/swaggerapi/",
}

// public handler for Kubernetes
func KubeEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	Log.Printf("------> KubeEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// check if uri pattern is accepted
	if ! IsSupportedPattern(r.RequestURI, kubePatterns){
		Log.Printf("Kube pattern not accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
		NoEndpointHandler(w, r)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	//Call Auth interceptor to authenticate with ccsapi
	status, node, namespace := auth.KubeAuth(r)
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

	//Handle request
	kubeHandler(w, r, node, namespace, req_id)

	Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
}

// private handler processing
func kubeHandler(w http.ResponseWriter, r *http.Request, redirect_host string,
	redirect_resource_id string, req_id string) {

	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false

	var err error = nil

	data, _ := httputil.DumpRequest(r, true)
	Log.Printf("Request dump of %d bytes:\n%s", len(data), string(data))

	body, _ := ioutil.ReadAll(r.Body)

	//***** Filter req/headers here before forwarding request to server *****

	if (httphelper.IsUpgradeHeader(r.Header)) {
		Log.Printf("@ Upgrade request detected\n")
		req_UPGRADE = true
	}

	maxRetries := 1
	backOffTimeout := 0

	var (resp *http.Response
		cc *httputil.ClientConn
	)
	for i:=0; i<maxRetries; i++ {
		resp, err, cc = redirect (r, body, redirect_host, redirect_resource_id,
			kubeRewriteUri, true /* override tls setting*/)
		if err == nil {
			break
		}
		Log.Printf("redirect retry=%d failed", i)
		if (i+1) < maxRetries {
			Log.Printf("will sleep secs=%d before retry", backOffTimeout)
			time.Sleep( time.Duration(backOffTimeout) * time.Second)
		}
	}
	if (err != nil) {
		Log.Printf("Error in redirection, will abort req_id=%s ... err=%v\n", req_id, err)
		return
	}

	//write out resp
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

	//TODO ***** Filter framework for Interception of commands before forwarding resp to client (1) *****

	proto := strings.ToUpper(httphelper.GetHeader(resp.Header, "Upgrade"))
	if (req_UPGRADE || resp_UPGRADE) && (proto != "TCP") {
		Log.Printf("Warning: will start hijack proxy loop although Upgrade proto %s is not TCP\n", proto)
	}

	if req_UPGRADE || resp_UPGRADE || resp_STREAM {
		//resp header is sent first thing on hijacked conn
		w.WriteHeader(resp.StatusCode)

		Log.Printf("starting tcp hijack proxy loop\n")
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

	_KUBE_CHUNKED_READ_ := true   // new feature flag

	if _KUBE_CHUNKED_READ_ {
		//new code to test
		//defer resp.Body.Close()   // causes this method to not return to caller IF closing while there is still data in Body!
		chunkedRWLoop(resp, w, req_id)
	}else {
		resp_body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			Log.Printf("Error: error in reading server response body\n")
			fmt.Fprint(w, "error in reading server response body\n")
			return
		}

		//TODO ***** Filter framework for Interception of commands before returning result to client (2) *****

		//Printout the response body
		bodystr := "Dump Body:\n"
		bodystr += httphelper.PrettyJson(resp_body)
		Log.Println(bodystr)
		/*
		if strings.ToLower(httphelper.GetHeader(resp.Header, "Content-Type")) == "application/json" {
			httphelper.PrintJson(resp_body)
		}else {
			Log.Printf("\n%s\n", string(resp_body))
		}
		*/
		//forward server response to calling client
		fmt.Fprintf(w, "%s", resp_body)
	}
	return
}

func kubeRewriteUri(reqUri string, namespace string) (redirectUri string){
	sl := strings.Split(reqUri, "/")
	next := false
	for i:=0; i < len(sl); i++{
		if next{
			redirectUri += namespace
			next = false
		}else{
			redirectUri += sl[i]
		}
		if sl[i] == "namespaces"{
			next = true
		}
		//if not done
		if i+1 < len(sl){
			redirectUri += "/"
		}
	}
	Log.Printf("kubeRewriteURI: '%s' --> '%s'\n", reqUri, redirectUri)
	return redirectUri
}
