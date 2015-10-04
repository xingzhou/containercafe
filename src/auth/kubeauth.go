package auth

// Kubernetes auth, and uri patterns

import (
	"net/http"
	"strings"
	"strconv"

	"conf"
)

func KubeAuth(r *http.Request) (status int, node string, namespace string) {
	status, host := getHost(r, "NoneContainer")
	if status == 200 {
		kubeMgr := injectKubePort( host.Mgr_host, conf.GetKubePort() ) 	// Kube master port is 6443
		Log.Printf("status=%d Mgr_host=%s namespace=%s", status, kubeMgr, host.Space_id)
		return status, kubeMgr, host.Space_id
	}
	Log.Printf("status=%d Mgr_host=\"\" namespace=\"\" ", status)
	return status, "", ""
}

func injectKubePort(host string, port int) string{
	//replace port by Kube's master port
	sl := strings.Split(host,":")
	return sl[0]+":"+strconv.Itoa(port)
}
