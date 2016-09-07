//kubernetes handler
//
package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"auth"
	"conf"
	"encoding/json"
	"httphelper"
	//	"reflect"
	"errors"

	"github.com/golang/glog"
)

// supported Kubernetes api uri prefix patterns
// these kube url patterns require namespaces:
var kubePrefixPatterns = []string{
	"/apis/",
	"/api/v1/namespaces/",
	"/api/v1/watch/namespaces/",
	"/api/v1/proxy/namespaces/",
	"/apis/",
	"/api/v1",
	"/apis/extensions",
	"/swaggerapi/",
}

// TODO
// There is a problem with /apis, because it's handled by docker endpoint pattern
//2016/04/18 15:58:51.066842 docker.go:75: ------> DockerEndpointHandler triggered, req_id=2, URI=/apis
//2016/04/18 15:58:51.066883 docker.go:79: Docker pattern not accepted, req_id=2, URI=/apis
//2016/04/18 15:58:51.066892 noendpoint.go:10: NoEndpointHandler triggered, URI=/apis, returning error 404
//2016/04/18 15:58:51.066903 docker.go:81: ------ Completed processing of request req_id=2

// these kube url patterns don't require namespaces
var kubeExactPatterns = []string{
	"/api",
	"/apis",
	"/version",
}

//called from init() of the package
func InitKubeHandler() {

}

// public handler for Kubernetes
func KubeEndpointHandler(w http.ResponseWriter, r *http.Request) {
	req_id := conf.GetReqId()
	glog.Info("Starting the Kube")
	glog.Infof("------> KubeEndpointHandler triggered, req_id=%s, URI=%s\n", req_id, r.RequestURI)

	// check if URI supported and requires auth.
	if IsExactPattern(r.RequestURI, kubeExactPatterns) {
		glog.Infof("Kube exact pattern accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
	} else if IsSupportedPattern(r.RequestURI, kubePrefixPatterns) {
		glog.Infof("Kube prefix pattern accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
	} else {
		glog.Infof("Kube pattern not accepted, req_id=%s, URI=%s", req_id, r.RequestURI)
		NoEndpointHandler(w, r)
		glog.Infof("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	glog.Infof("This is a AUTH Kube supported pattern %+v", r.RequestURI)
	// body, _ := ioutil.ReadAll(r.Body)
	// Log.Printf("**** %+v", r)
	// Log.Printf("**** This is a request body: %+v", body)

	// read the credentials from the local file first
	var creds auth.Creds
	creds = auth.FileAuth(r) // So creds should now hold info FOR THAT space_id.
	if creds.Status == 200 {
		glog.Infof("Authentication from FILE succeeded for req_id=%s status=%d", req_id, creds.Status)
		// Log.Printf("**** Creds %+v", creds)
	} else {
		glog.Errorf("Authentication from FILE failed for req_id=%s status=%d", req_id, creds.Status)
		if creds.Status == 401 {
			NotAuthorizedHandler(w, r)
		} else {
			ErrorHandler(w, r, creds.Status)
		}
		glog.Infof("------ Completed processing of request req_id=%s\n", req_id)
		return
	}

	// validate the creds
	if creds.Node == "" || creds.Space_id == "" {
		glog.Errorf("Missing data. Host = %v, Space_id = %v", creds.Node, creds.Space_id)
		ErrorHandlerWithMsg(w, r, 404, "Incomplete data received from authentication component")
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

	glog.Infof("Assigning proper Kubernetes port. Old target: %v, New target: %v", creds.Node, redirectTarget)

	// get user certificates from the CCSAPI server
	status, certs := auth.GetCert(r, creds)
	//status, certs := auth.GetCert(r)
	if status != 200 {
		glog.Errorf("Obtaining user certs failed for req_id=%s status=%d", req_id, status)
		ErrorHandler(w, r, creds.Status)
	}
	glog.Infof("Obtaining user certs successful for req_id=%s status=%d", req_id, status)

	// convert the Bluemix space id to namespace
	namespace := auth.GetNamespace(creds.Space_id)
	kubeHandler(w, r, redirectTarget, namespace, req_id, []byte(certs.User_cert), []byte(certs.User_key))
	glog.Infof("------ Completed processing of request req_id=%s\n", req_id)
}

// private handler processing
func kubeHandler(w http.ResponseWriter, r *http.Request, redirect_host string,
	namespace string, req_id string, cert []byte, key []byte) {

	req_UPGRADE := false
	resp_UPGRADE := false
	resp_STREAM := false

	var err error = nil

	data, _ := httputil.DumpRequest(r, true)
	glog.Infof("Request dump of %d bytes:\n%s", len(data), string(data))
	glog.Infof("Redirect host %v\n", redirect_host)
	//body, _ := ioutil.ReadAll(r.Body)

	// sometimes body needs to be modify to add custom labels, annotations
	body, err := kubeUpdateBody(r, namespace)
	if err != nil {
		glog.Errorf("Error %v", err.Error())
		ErrorHandlerWithMsg(w, r, 500, "Error updating Kube body: "+err.Error())
	}

	//***** Filter req/headers here before forwarding request to server *****

	if httphelper.IsUpgradeHeader(r.Header) {
		glog.Infof("@ Upgrade request detected\n")
		req_UPGRADE = true
	}

	maxRetries := 1
	backOffTimeout := 0

	var (
		resp *http.Response
		cc   *httputil.ClientConn
	)
	for i := 0; i < maxRetries; i++ {
		// resp, err, cc = redirect_random (r, body, redirect_host, redirect_resource_id,
		resp, err, cc = redirect_with_cert(r, body, redirect_host, namespace,
			kubeRewriteUri, false, cert, key /* override tls setting*/)
		// kubeRewriteUri, true /* override tls setting*/) TODO MS
		if err == nil {
			break
		}
		glog.Warningf("redirect retry=%d failed", i)
		if (i + 1) < maxRetries {
			glog.Warningf("will sleep secs=%d before retry", backOffTimeout)
			time.Sleep(time.Duration(backOffTimeout) * time.Second)
		}
	}
	if err != nil {
		glog.Errorf("Error in redirection, will abort req_id=%s ... err=%v\n", req_id, err)
		ErrorHandlerWithMsg(w, r, 500, "Internal communication error. Check if the redirected host is active")
		return
	}

	//write out resp
	glog.Infof("<------ req_id=%s\n", req_id)
	//data2, _ := httputil.DumpResponse(resp, true)
	//fmt.Printf("Response dump of %d bytes:\n", len(data2))
	//fmt.Printf("%s\n", string(data2))

	glog.Infof("Resp Status: %s\n", resp.Status)
	glog.Info(httphelper.DumpHeader(resp.Header))

	httphelper.CopyHeader(w.Header(), resp.Header)

	if httphelper.IsUpgradeHeader(resp.Header) {
		glog.Infof("@ Upgrade response detected\n")
		resp_UPGRADE = true
	}
	if httphelper.IsStreamHeader(resp.Header) {
		glog.Infof("@ application/octet-stream detected\n")
		resp_STREAM = true
	}

	//TODO ***** Filter framework for Interception of commands before forwarding resp to client (1) *****

	proto := strings.ToUpper(httphelper.GetHeader(resp.Header, "Upgrade"))
	if (req_UPGRADE || resp_UPGRADE) && (proto != "TCP") {
		glog.Warningf("Warning: will start hijack proxy loop although Upgrade proto %s is not TCP\n", proto)
	}

	if req_UPGRADE || resp_UPGRADE || resp_STREAM {
		//resp header is sent first thing on hijacked conn
		w.WriteHeader(resp.StatusCode)

		glog.Infof("starting tcp hijack proxy loop\n")
		httphelper.InitProxyHijack(w, cc, req_id, "TCP") // TCP is the only supported proto now
		return
	}
	//If no hijacking, forward full response to client
	w.WriteHeader(resp.StatusCode)

	if resp.Body == nil {
		glog.Infof("\n")
		fmt.Fprintf(w, "\n")
		return
	}

	_KUBE_CHUNKED_READ_ := true // new feature flag

	if _KUBE_CHUNKED_READ_ {
		//new code to test
		//defer resp.Body.Close()   // causes this method to not return to caller IF closing while there is still data in Body!
		chunkedRWLoop(resp, w, req_id)
	} else {
		resp_body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			glog.Errorf("Error: error in reading server response body\n")
			fmt.Fprint(w, "error in reading server response body\n")
			return
		}

		//TODO ***** Filter framework for Interception of commands before returning result to client (2) *****

		//Printout the response body
		bodystr := "Dump Body:\n"
		bodystr += httphelper.PrettyJson(resp_body)
		glog.Info(bodystr)
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

func kubeUpdateBody(r *http.Request, namespace string) (body []byte, err error) {
	body, _ = ioutil.ReadAll(r.Body)
	if r.Method != "POST" {
		// return the original body
		return body, nil
	}
	// convert the body to string
	bodystr := httphelper.PrettyJson(body)
	glog.Infof("Original JSON: %s", bodystr)

	// the request to create pod looks as follow:
	//	 {
	//	  	"kind":"Pod",
	//	  	"apiVersion:"v1",
	//	  	"metadata":{
	//	  		"name":"testtt1",
	//	  		"namespace":"s21f85bc8-5a1a-403a-8a82-cdb757defd72-default",
	//	  		"annotations":{
	//	  			containers-annotations.alpha.kubernetes.io: "{ \"com.ibm.radiant.tenant.0\": \"stest1-default\",  \"OriginalName\": \"kube-web-server\" }"
	//
	// and the one to create depoloyment (group):
	//  	"kind": "Deployment",
	//		"apiVersion": "extensions/v1beta1",
	//		"metadata": {
	//			"name": "k3",
	//			"creationTimestamp": null,
	//			"labels": {
	//				"run": "k3"
	//			}
	//		},
	//		"spec": {
	//			"replicas": 1,
	//			"selector": {
	//				"matchLabels": {
	//					"run": "k3"
	//				}
	//			},
	//			"template": {
	//				"metadata": {
	//					"creationTimestamp": null,
	//					"annotations":{
	//	  					containers-annotations.alpha.kubernetes.io: "{ \"com.ibm.radiant.tenant.0\": \"stest1-default\",  \"OriginalName\": \"kube-web-server\" }"
	//					"labels": {
	//						"run": "k3"
	//					}
	//				},}}}

	// get the label names
	auth_label := conf.GetSwarmAuthLabel()
	annot_label := conf.GetAnnotationExtLabel()

	// get the body of the request
	data := map[string]interface{}{}
	json.Unmarshal(body, &data)

	// TODO the code below has to be revisited and simplified
	// there should be a common method for injecting the annotations
	// I just could get this working yet....
	//		fmt.Println(reflect.TypeOf(data))
	//		meta := inject_annotation(data["metadata"])

	// get request type
	kind := data["kind"]

	var metam map[string]interface{}
	if kind == "Pod" {
		meta := data["metadata"]
		// convert the interface{} to map
		metam = meta.(map[string]interface{})
	} else if kind == "Deployment" || kind == "ReplicaSet" || kind == "ReplicationController" || kind == "Job" {
		spec := data["spec"]
		specm := spec.(map[string]interface{})
		templ := specm["template"]
		templm := templ.(map[string]interface{})
		meta := templm["metadata"]
		metam = meta.(map[string]interface{})
	}

	annot := metam["annotations"]
	//Log.Printf("Original Annotations: %+v", annot)
	var annotm map[string]interface{}
	// annotations might not be provided in the yaml file:
	if annot == nil {
		annotm = make(map[string]interface{})
		metam["annotations"] = annotm
	} else {
		// convert the existing annotation interface{} to map
		annotm = annot.(map[string]interface{})
	}
	if annotm[annot_label] == "" || annotm[annot_label] == nil {
		glog.Infof("Annotation label does not exist")
	} else {
		glog.Infof("Annotation label %v already exists: %v", annot_label, annotm[annot_label])
		err = errors.New("Illegal usage of label ")
		return nil, err
	}

	new_value := "{ \"" + auth_label + "\": \"" + namespace + "\" }"
	annotm[annot_label] = new_value
	glog.Infof("Injecting annotation: %v with value: %v for kind: %v", annot_label, new_value, kind)

	b, err := json.Marshal(data)
	if err != nil {
		glog.Errorf("Error marshaling the updated json: %v ", err)
		return nil, err
	}
	bodystr = httphelper.PrettyJson(b)
	glog.Info("Updated JSON: %s", bodystr)
	return b, nil
}

func kubeRewriteUri(reqUri string, namespace string) (redirectUri string) {
	sl := strings.Split(reqUri, "/")
	next := false
	for i := 0; i < len(sl); i++ {
		if next {
			redirectUri += namespace
			next = false
		} else {
			redirectUri += sl[i]
		}
		if sl[i] == "namespaces" {
			next = true
		}
		//if not done
		if i+1 < len(sl) {
			redirectUri += "/"
		}
	}
	glog.Infof("kubeRewriteURI: '%s' --> '%s'\n", reqUri, redirectUri)
	return redirectUri
}
