package auth

import (
	"net/http"
	"encoding/json"
	"io"
	"os"

	"conf"  		// my conf package
)

// Use a file as authentication credentials store (mainly for trusted test SWARM tenants)
// BlueMix space id is passed in request header as X-Auth-Project-Id header and is used as search key into the file
func FileAuth(r *http.Request) (creds Creds) {
	creds.Status = 404
	proxy_auth_header := r.Header.Get("X-Auth-Proxy")
	if proxy_auth_header != "TOKEN" && proxy_auth_header != "Token" && proxy_auth_header != "token" {
		return
	}
	space_id := r.Header.Get("X-Auth-Project-Id")

	fname := conf.GetStubAuthFile()
	fp, err := os.Open(fname)
	if err != nil{
		Log.Println(err)
		return
	}

	dec := json.NewDecoder(fp)

	for {
		var c Creds
		if err := dec.Decode(&c); err == io.EOF {
			break
		} else if err != nil {
			Log.Println(err)
			break
		}

		if space_id == c.Space_id {
			creds = c
			//creds.Status = 200  //return the status that is in the auth conf file
			//Set Swarm Authorization header
			r.Header.Set("X-Auth-Token", space_id)
			return
		}
	}

	//tenant not found in credentials file
	creds.Status = 401
	return
}
