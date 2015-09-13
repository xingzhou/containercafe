package auth

// Kubernetes auth, and uri patterns

import (
	"net/http"
	"strings"
	"strconv"
	"log"

	"conf"
)

func KubeAuth(r *http.Request) (ok bool, node string, namespace string) {
	ok, host := getHost(r, "NoneContainer")
	if ok {
		kubeMgr := injectKubePort( host.Mgr_host, conf.GetKubePort() ) 	// Kube master port is 6443
		log.Printf("ok=true, Mgr_host=%s, namespace=%s", kubeMgr, host.Space_id)
		return ok, kubeMgr, host.Space_id
	}
	log.Println("ok=false, Mgr_host=\"\", namespace=\"\" ")
	return false, "", ""
}

func injectKubePort(host string, port int) string{
	//replace port by Kube's master port
	sl := strings.Split(host,":")
	return sl[0]+":"+strconv.Itoa(port)
}
