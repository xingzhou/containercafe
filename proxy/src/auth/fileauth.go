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
	// let's start with 400 Bad Request 
	creds.Status = 400
	if (r.TLS == nil || len(r.TLS.PeerCertificates) < 1) {
		Log.Printf("**** Error, request missing client TLS certificate")
		return
	}

	Log.Printf("**** RequestTLS %+v", r.TLS)
	Log.Printf("**** Length TLS %v", len(r.TLS.PeerCertificates))
	Log.Printf("**** Client TLS %v", r.TLS.PeerCertificates)
	
	cn := ""
	//var errlist []error
	for _, cert := range r.TLS.PeerCertificates {
		
		Log.Printf("**** CN from CERT: %v", cert.Subject.CommonName)
		cn = cert.Subject.CommonName
		// cert could be CA:
		if (len(cn) != 0 && cn != "containers-api-dev.stage1.ng.bluemix.net") {
			break
		}
		
		Log.Printf("**** Subject from CERT: %+v", cert.Subject)
	}
	if cn == "" {
		return
	}
	
	Apikey := cn
	//	chains, err := cert.Verify(a.opts)
	//	if err != nil {
	//		errlist = append(errlist, err)
	//		continue
	//	}

//	for _, chain := range chains {
//		user, ok, err := a.user.User(chain)
	//	if err != nil {
//			errlist = append(errlist, err)
//			continue
//		}

//		if ok {
//			return user, ok, err
//		}
//	}
	
	
	
    //Log.Printf("**** Request TLSUnique %+v", r.TLS.TLSUnique)
    
	creds.Status = 404
	//  swarm-auth now uses 'X-Auth-TenantId' instead of 'X-Auth-Project-Id'
	// space_id := r.Header.Get("X-Auth-Project-Id")
	//space_id := r.Header.Get(conf.GetSwarmAuthHeader())
	//Apikey := r.Header.Get("X-Tls-Client-Dn")
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
		// skip if this is not swarm shard nor radiant endpoint
		if (c.Endpoint_type != "radiant" && c.Swarm_shard != true) {
			continue
		}
		
		//Log.Printf("****** API_ID=%v, creds.Apikey=%v", api_id, c.Apikey)
		//credkey := "/CN="+c.Apikey
		if Apikey == c.Apikey || Apikey == ("/CN=" + c.Apikey) {
			space_id := c.Space_id
			creds = c
			//creds.Status = 200  //return the status that is in the auth conf file
			//Set Swarm Authorization header
			namespace := GetNamespace(space_id)
			r.Header.Set(conf.GetSwarmAuthHeader(), namespace)
			return
		}
	}
    	Log.Printf("Tenant API key %s not found in %s file", Apikey, fname)
	//tenant not found in credentials file
	creds.Status = 401
	return
}
