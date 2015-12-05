package auth

import (
	"net/http"
)

var validSpaces = []Creds {
	{200, "10.140.28.132:2379", "", "", true, "85cdc7e0-32d8-4552-9bae-907c3f1d98d9", "swarm", "c3d87893a5b7f56991fd328f655f25cce286591c3ce4a558","9013217d-0abf-40fe-bd35-bb625066408c", "924fc412d1004528b90007e898aeb0d8"},
}

// authentication stub for trusted SWARM tenants
// BlueMix space id is passed in request as X-Auth-Project-Id header
func StubAuth(r *http.Request) (creds Creds) {
	proxy_auth_header := r.Header.Get("X-Auth-Proxy")
	if proxy_auth_header != "TOKEN" && proxy_auth_header != "Token" && proxy_auth_header != "token" {
		creds.Status = 404
		return
	}

	space_id := r.Header.Get("X-Auth-Project-Id")
	for _, v := range validSpaces {
		if space_id == v.Space_id {
			creds = v
			creds.Status = 200

			//Set Swarm Authorization header
			r.Header.Set("X-Auth-Token", space_id)

			creds.Docker_id = "" // --> no rewrite of docker_id  (rewrite is not needed with Swarm, needed in nova-docker case only)
			return
		}
	}
	creds.Status = 401
	return
}
