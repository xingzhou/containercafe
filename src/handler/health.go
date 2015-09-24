package handler

import (
	"net/http"
	"log"
	"fmt"
	"strings"

	"conf"
)

// supported health api uri prefix patterns
var healthPatterns = []string {
	"/hjproxy/health",
	"/hjproxy/stats",
	"/hjproxy/_ping_notls", // _ping with no tls override regardless of hjproxy configuration
	"/hjproxy/_ping",  		//  /hjproxy/_ping/host   or   /hjproxy/_ping/host/port
}

func HealthEndpointHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HealthEndpointHandler triggered, URI=%s\n", r.RequestURI)
	p := GetUriPattern(r.RequestURI, healthPatterns)
	switch p{
	case healthPatterns[0]:
		v := conf.GetVerStr()
		fmt.Fprintf(w,"hjproxy up\n%s\n", v)
		log.Println("hjproxy up", v)
		break
	case healthPatterns[1]:
		v := conf.GetVerStr()
		n := conf.GetNumServedRequests()
		fmt.Fprintf(w,"hjproxy %s\n", v)
		fmt.Fprintf(w,"This instance served %d requests\n", n)
		log.Println("hjproxy", v)
		log.Printf("This instance served %d requests\n", n)
		break
	case healthPatterns[2]:
		ping(w, r, true)
		break
	case healthPatterns[3]:
		ping(w, r, false)
		break
	default:
		log.Printf("Health pattern not accepted, URI=%s", r.RequestURI)
		NoEndpointHandler(w, r)
	}
}

func ping(w http.ResponseWriter, r *http.Request, tls_override bool) {
	var status int

	sl := strings.Split(r.RequestURI, "/")
	server := sl[3]
	port := conf.GetDockerPort()
	if len(sl)>4 {
		port = sl[4]
	}
	redirect_host := server + ":" + port

	resp, err, _ := redirect (r, nil /*body*/, redirect_host, "" /*resource_id*/, pingRewriteUri, tls_override)
	if (err != nil) {
		log.Printf("Error in redirection of _ping to %s ... err=%v\n", redirect_host, err)
		status = 500
	}else {
		status = resp.StatusCode
	}

	if status == 200 {
		log.Printf("_ping success to host=%s\n", redirect_host)
		fmt.Fprintf(w,"_ping success to host=%s\n", redirect_host)  //returns status 200
	}else{
		ErrorHandler(w, r, status)
	}
}

func pingRewriteUri(reqUri string, resource_id string) (newReqUri string){
	newReqUri = conf.GetDockerApiVer()+"/_ping"
	log.Printf("pingRewriteURI: '%s' --> '%s'\n", reqUri, newReqUri)
	return newReqUri
}
