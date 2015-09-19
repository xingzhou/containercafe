package handler

import (
	"net/http"
	"log"
	"fmt"

	"conf"
)

// supported health api uri prefix patterns
var healthPatterns = []string {
	"/hjproxy/health",
	"/hjproxy/stats",
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
	default:
		log.Printf("Health pattern not accepted, URI=%s", r.RequestURI)
		NoEndpointHandler(w, r)
	}
}

