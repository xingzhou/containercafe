package auth

// Groups auth

import (
	"net/http"
)

func GroupsAuth(r *http.Request) (creds Creds) {
	var host GetHostResp
	creds.Status, host = getHost(r, "NoneContainer")

	creds.Node = host.Groups_host
	creds.Reg_namespace = host.Namespace
	creds.Apikey = host.Apikey
	creds.Space_id = GetNamespace(host.Space_id)
	creds.Orguuid = host.Orguuid
	creds.Userid = host.Userid
	creds.Tls_override = ! host.Swarm_tls

	Log.Printf("status=%d Groups_host=%s Space_id=%s", creds.Status, creds.Node, creds.Space_id)
	return
}

