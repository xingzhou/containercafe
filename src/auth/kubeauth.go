package auth

// Kubernetes auth, and uri patterns

import (
	"net/http"
	"strings"
	"strconv"

	//"conf"
)

//KubeAuth uses only the following fields of Creds[]: Status, Node, Space_id
func KubeAuth(r *http.Request) (creds Creds) {
	var host GetHostResp
	Log.Printf("Processing KubeAuth")
	creds.Status, host = getHost(r, "NoneContainer")
	if creds.Status == 200 {
		//kubeMgr := injectKubePort( host.Mgr_host, conf.GetKubePort() ) 	// Kube master port is 6443
		creds.Node = host.Mgr_host // radiant manager
		creds.Space_id = GetNamespace(host.Space_id)
	}
	Log.Printf("status=%d Mgr_host=%s namespace=%s", creds.Status, creds.Node, creds.Space_id)
	return
}

func injectKubePort(host string, port int) string{
	//replace port by Kube's master port
	sl := strings.Split(host,":")
	return sl[0]+":"+strconv.Itoa(port)
}
