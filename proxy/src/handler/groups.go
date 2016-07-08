// Groups API handler
//
package handler

import (
    "fmt"
    "net/http"
	"net/http/httputil"
	"io/ioutil"
	"time"
	"strings"

	"httphelper"  	//my httphelper package
	"auth"  		// my auth package
	"conf"  		// my conf package
	//"limit"  		//my limits package
)

// supported docker api uri prefixes
var groupsPatterns = []string {
	"/groups",
	"/groups/",
}

// Router based on uri patterns wih simple expressions
var groupsRouter *Router

//called from init() of the handler package, before any requests are handled
func InitGroupsHandler(){
	//define routes for api endpoints
	groupsRoutes := []Route{
		NewRoute("*", "*", groupsHandler),  //wildcard for forwarding everything else
	}
	groupsRouter = NewRouter(groupsRoutes)
}

// http proxy forwarding with hijack support
// handler for docker/swarm
func GroupsEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	Log.Printf("------> GroupsEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// check if uri pattern is accepted
	if ! IsSupportedPattern(r.RequestURI, groupsPatterns){
		Log.Printf("Groups pattern not accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
		NoEndpointHandler(w, r)
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	data, _ := httputil.DumpRequest(r, true)
	Log.Printf("Request dump req_id=%s req_length=%d:\n%s", req_id, len(data), string(data))

	var creds auth.Creds

	// workaround defective sharding in dev-mon
	creds = auth.FileAuth(r)
	if creds.Status == 200 {
		Log.Printf("Authentication from FILE succeeded for req_id=%s status=%d", req_id, creds.Status)
	}else {
		creds = auth.GroupsAuth(r)
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

	body, _ := ioutil.ReadAll(r.Body)

	/*
	//Call conn limiting interceptor(s) pre-processing,
	if !limit.OpenConn(creds.Container, conf.GetMaxContainerConn()) {
		Log.Printf("Max conn limit reached for container...aborting request")
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
	if !limit.OpenConn(creds.Node, conf.GetMaxNodeConn()) {
		Log.Printf("Max conn limit reached for host node...aborting request")
		Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
		return
	}
	*/

	//Handle request
	groupsRouter.DoRoute(w, r, body, creds, req_id)

	/*
	//Call conn limiting interceptor(s) post-processing, to decrement conn count(s)
	limit.CloseConn(creds.Container, conf.GetMaxContainerConn())
	limit.CloseConn(creds.Node, conf.GetMaxNodeConn())
	*/

	Log.Printf("------ Completed processing of request req_id=%s\n", req_id)
}

/////////////////
// route handlers
/////////////////

func groupsNotSupported(w http.ResponseWriter, r *http.Request, body []byte, creds auth.Creds, vars map[string]string, req_id string){
	Log.Printf("Groups pattern not accepted, URI=%s", r.RequestURI)
	NoEndpointHandler(w, r)
}

// default route handler
func groupsHandler(w http.ResponseWriter, r *http.Request, body []byte, creds auth.Creds, vars map[string]string, req_id string) {
	redirect_host := creds.Node
	redirect_resource_id := ""
	tls_override := creds.Tls_override
	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false
	var err error = nil

	//***** Filter req/headers here before forwarding request to server *****
	if (httphelper.IsUpgradeHeader(r.Header)) {
		Log.Printf("@ Upgrade request detected\n")
		req_UPGRADE = true
	}

	maxRetries := 1 //conf.GetMaxRetries()
	backOffTimeout := 0 //conf.GetBackOffTimeout()

	Log.Printf("Redirecting to host=%s req_id=%s", redirect_host, req_id)
	var (resp *http.Response
		cc *httputil.ClientConn
	)
	for i:=0; i<maxRetries; i++ {
		resp, err, cc = redirect_random(r, body, redirect_host, redirect_resource_id,
			groupsRewriteUri, tls_override)
		if err == nil {
			break
		}
		Log.Printf("redirect failed retry=%d req_id=%s err=%s", i, req_id, err)
		if (i+1) < maxRetries {
			Log.Printf("will sleep before retry secs=%d req_id=%s", backOffTimeout, req_id)
			time.Sleep( time.Duration(backOffTimeout) * time.Second)
		}
	}
	if (err != nil) {
		Log.Printf("Error in redirection, will abort req_id=%s err=%v\n", req_id, err)
		return
	}

	Log.Printf("<------ req_id=%s\n", req_id)
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

	//***** Filter framework for Interception of commands before forwarding resp to client (1) *****

	proto := strings.ToUpper(httphelper.GetHeader(resp.Header, "Upgrade"))
	if (req_UPGRADE || resp_UPGRADE) && (proto != "TCP") {
		Log.Printf("Warning: will start hijack proxy loop although Upgrade proto %s is not TCP\n", proto)
	}

	if req_UPGRADE || resp_UPGRADE || resp_STREAM {
		//resp header is sent first thing on hijacked conn
		w.WriteHeader(resp.StatusCode)

		Log.Printf("starting tcp hijack proxy loop req_id=%s", req_id)
		httphelper.InitProxyHijack(w, cc, req_id, "TCP") // TCP is the only supported proto now
		return
	}
	//If no conn upgrade, forward full response to client
	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		Log.Printf("\n")
		fmt.Fprintf(w, "\n")
		return
	}

	_GROUPS_CHUNKED_READ_ := true   // feature flag

	if _GROUPS_CHUNKED_READ_ {
		//defer resp.Body.Close()   // causes this method to not return to caller IF closing while there is still data in Body!
		chunkedRWLoop(resp, w, req_id)

	}else {
		resp_body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			Log.Printf("Error: error in reading server response body")
			fmt.Fprint(w, "error in reading server response body\n")
			return
		}

		//Printout the response body
		bodystr := "Dump Body:\n"
		bodystr += httphelper.PrettyJson(resp_body)
		Log.Println(bodystr)

		//forward server response to calling client
		fmt.Fprintf(w, "%s", resp_body)
	}
	return
}

func groupsRewriteUri(reqUri string, redirect_resource_id string)(redirectUri string){
	// no URI transformation
	redirectUri = reqUri
	Log.Printf("groupsRewriteURI: '%s' --> '%s'\n", reqUri, redirectUri)
	return redirectUri
}
