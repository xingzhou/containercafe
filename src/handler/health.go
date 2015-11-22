package handler

import (
	"net/http"
	"fmt"
	"strings"

	"conf"
)

// supported health api uri prefix patterns
var healthPatterns = []string {
	"/hjproxy/health",
	"/hjproxy/stats",
	"/hjproxy/_ping_notls", // _ping with no tls override regardless of hjproxy configuration
	"/hjproxy/_ping",  		//  /hjproxy/_ping/host  or  /hjproxy/_ping/host/port  8089 is used if port not specified
}

func HealthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	Log.Printf("HealthEndpointHandler triggered, URI=%s", r.RequestURI)
	p := GetUriPattern(r.RequestURI, healthPatterns)
	switch p{
	case healthPatterns[0]:
		v := conf.GetVerStr()
		fmt.Fprintf(w,"hjproxy up\n%s\n", v)
		Log.Print("hjproxy up ", v)
		break
	case healthPatterns[1]:
		v := conf.GetVerStr()
		n := conf.GetNumServedRequests()
		fmt.Fprintf(w,"hjproxy %s\n", v)
		fmt.Fprintf(w,"This instance served %d non-admin requests\n", n)
		Log.Print("hjproxy ", v)
		Log.Printf("This instance served %d non-admin requests", n)
		break
	case healthPatterns[2]:
		ping(w, r, true)
		break
	case healthPatterns[3]:
		ping(w, r, false)
		break
	default:
		Log.Printf("Health pattern not accepted, URI=%s", r.RequestURI)
		NoEndpointHandler(w, r)
	}
}

func ping(w http.ResponseWriter, r *http.Request, tls_override bool) {
	var status int

	sl := strings.Split(r.RequestURI, "/")
	if len(sl)<=3 {
		Log.Print("Error in _ping ... no host specified")
		NoEndpointHandler(w, r)
		return
	}
	server := sl[3]
	port := conf.GetDockerPort()
	if len(sl)>4 {
		port = sl[4]
	}
	redirect_host := server + ":" + port

	resp, err, _ := redirect (r, nil /*body*/, redirect_host, "" /*resource_id*/, pingRewriteUri, tls_override)
	if (err != nil) {
		Log.Printf("Error in redirection of _ping to %s ... err=%v", redirect_host, err)
		status = 500
	}else {
		status = resp.StatusCode
	}

	if status == 200 {
		Log.Printf("_ping success to host=%s", redirect_host)
		fmt.Fprintf(w,"_ping success to host=%s\n", redirect_host)  //returns status 200
	}else{
		//ErrorHandler(w, r, status)
		msg := fmt.Sprintf("_ping failed to host=%s", redirect_host)
		ErrorHandlerWithMsg(w, r, status, msg)
	}
}

func pingRewriteUri(reqUri string, resource_id string) (newReqUri string){
	newReqUri = conf.GetDockerApiVer()+"/_ping"
	Log.Printf("pingRewriteURI: '%s' --> '%s'", reqUri, newReqUri)
	return newReqUri
}
