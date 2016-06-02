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
	"strconv"

	"httphelper"
	"conf"
	"auth"
)

// supported Kubernetes api uri prefix patterns
// these kube url patterns require namespaces:
var kubePrefixPatterns = []string {
	"/api/v1/namespaces/",
	"/api/v1/watch/namespaces/",
	"/api/v1/proxy/namespaces/",
	"/apis/",
	"/swaggerapi/",
}

// TODO 
// There is a problem with /apis, because it's handled by docker endpoint pattern
//2016/04/18 15:58:51.066842 docker.go:75: ------> DockerEndpointHandler triggered, req_id=2, URI=/apis
//2016/04/18 15:58:51.066883 docker.go:79: Docker pattern not accepted, req_id=2, URI=/apis
//2016/04/18 15:58:51.066892 noendpoint.go:10: NoEndpointHandler triggered, URI=/apis, returning error 404
//2016/04/18 15:58:51.066903 docker.go:81: ------ Completed processing of request req_id=2

// these kube url patterns don't require namespaces
var kubeExactPatterns = []string {
	"/api",
	"/apis",
	"/version",
}


//called from init() of the package
func InitKubeHandler(){

}

// public handler for Kubernetes
func KubeEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	Log.Printf("Starting the Kube")
	Log.Printf("------> KubeEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// check if URI supported and requires auth.
	if IsExactPattern(r.RequestURI, kubeExactPatterns){
		Log.Printf("Kube exact pattern accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
	} else if IsSupportedPattern(r.RequestURI, kubePrefixPatterns) {
		Log.Printf("Kube prefix pattern accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
	} else {
		Log.Printf("Kube pattern not accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
		NoEndpointHandler(w, r)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	Log.Printf("This is a AUTH Kube supported pattern %+v", r.RequestURI)
	// body, _ := ioutil.ReadAll(r.Body)
	// Log.Printf("**** %+v", r)
	// Log.Printf("**** This is a request body: %+v", body)
	
	// read the credentials from the local file first
	var creds auth.Creds
	creds = auth.FileAuth(r) // So creds should now hold info FOR THAT space_id. 
	if creds.Status == 200 {
		Log.Printf("Authentication from FILE succeeded for req_id=%s status=%d", req_id, creds.Status)
		Log.Printf("Will not execute CCSAPI auth")
		// Log.Printf("**** Creds %+v", creds)
	} else {
		Log.Printf("Authentication from FILE failed for req_id=%s status=%d", req_id, creds.Status)
		Log.Printf("Excuting CCSAPI auth")
		
		creds = auth.KubeAuth(r)
		// Log.Printf("***** Creds: %+v", creds)
	
		if creds.Status == 200 {
			Log.Printf("CCSAPI Authentication succeeded for req_id=%s status=%d", req_id, creds.Status)
		} else {
			Log.Printf("CCAPI Auth failed to process req_id=%s\n", req_id)
			if creds.Status == 401 {
				NotAuthorizedHandler(w, r)
			} else {
				ErrorHandler(w, r, creds.Status)
			}
			Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
			return
		}	
		Log.Printf("CCSAPI Authentication succeeded for req_id=%s status=%d", req_id, creds.Status)
	}

	// validate the creds
	if creds.Node == "" || creds.Space_id == "" {
		Log.Printf("Missing data. Host = %v, Space_id = %v", creds.Node, creds.Space_id)
		ErrorHandlerWithMsg(w, r, 404, "Incomplete data received from CCSAPI server")
		return
	}
	
	// assigning a proper port for Kubernentes
	// the target might or might not contain 'http://', strip it
	redirectTarget := creds.Node
	sp := strings.Split(creds.Node, ":")
	if sp[0] == "http" || sp[0] == "https" {
		redirectTarget = sp[1] + ":" + strconv.Itoa(conf.GetKubePort())
		// strip out the '//' from http://
		redirectTarget = redirectTarget[2:len(redirectTarget)]
	} else {
		redirectTarget = sp[0] + ":" + strconv.Itoa(conf.GetKubePort())
	}
	
	Log.Printf("Assigning proper Kubernetes port. Old target: %v, New target: %v", creds.Node, redirectTarget)
	
	// TODO for now skip the StubAuth, not needed
	if (false) {
	//Call Auth interceptor to authenticate with ccsapi
	creds = auth.StubAuth(r)
	Log.Printf("*** Creds %+v", creds)
	if creds.Status == 200 {
		Log.Printf("Stub Authentication succeeded for req_id=%s status=%d", req_id, creds.Status)
	}else {
		creds = auth.KubeAuth(r)
		if creds.Status != 200 {
			Log.Printf("Authentication failed for req_id=%s status=%d", req_id, creds.Status)
			if creds.Status == 401 {
				NotAuthorizedHandler(w, r)
			}else {
				ErrorHandler(w, r, creds.Status)
			}
			Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
			return
		}
		Log.Printf("Authentication succeeded for req_id=%s status=%d", req_id, creds.Status)
	}
    }

	// get user certificates from the CCSAPI server
	 status, certs := auth.GetCert(r, creds)
	 //status, certs := auth.GetCert(r)
	 if status != 200 {
	 	Log.Printf("Obtaining user certs failed for req_id=%s status=%d", req_id, status)
			ErrorHandler(w, r, creds.Status)
	 }
	 Log.Printf("Obtaining user certs successful for req_id=%s status=%d", req_id, status)
	
	// convert the Bluemix space id to namespace
	namespace := auth.GetNamespace(creds.Space_id)
	kubeHandler(w, r, redirectTarget, namespace, req_id, []byte(certs.User_cert), []byte(certs.User_key))
	Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
}

// private handler processing
func kubeHandler(w http.ResponseWriter, r *http.Request, redirect_host string,
	redirect_resource_id string, req_id string, cert []byte, key []byte) {

	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false

	var err error = nil

	data, _ := httputil.DumpRequest(r, true)
	Log.Printf("Request dump of %d bytes:\n%s", len(data), string(data))
	Log.Printf("Redirect host %v\n", redirect_host)
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
		// resp, err, cc = redirect_random (r, body, redirect_host, redirect_resource_id,
		resp, err, cc = redirect_with_cert(r, body, redirect_host, redirect_resource_id,
			kubeRewriteUri, false, cert, key /* override tls setting*/)
			// kubeRewriteUri, true /* override tls setting*/) TODO MS
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
		ErrorHandlerWithMsg(w, r, 500, "Internal communication error. Check if the redirected host is active")
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
