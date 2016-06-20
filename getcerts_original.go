package auth

import (
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"httphelper"  	// my httphelper
	"conf"  		// my conf package
)


// '/v3/tlskey' CCSAPI call response
type GetCertResp struct {
	User_cert  	string
	User_key   	string
	Server_cert 	string
	Ca_cert		string
}


// forward r header only without body to ccsapi auth endpoint
// return ok=true if resp status is 200, otherwise ok=false
func GetCert(r *http.Request) (status int, certs GetCertResp){
	new_uri := "http://"+conf.GetCcsapiHost()+"/v3/tlskey"

	if ! AuthHeadersExist(r.Header){
		Log.Println("Auth headers missing. Will NOT invoke CCSAPI to authenticate.")
		status = 500
		return
	}

	req, _ := http.NewRequest("GET", new_uri, nil)
	httphelper.CopyHeader(req.Header, r.Header)  //req.Header = r.Header
	req.URL.Host = conf.GetCcsapiHost()
	
	client := &http.Client{
		CheckRedirect: nil,
	}
	Log.Println("will invoke CCSAPI to authenticate and get certs...")

	resp, err := client.Do(req)
	if (err != nil) {
		Log.Printf("GetCert: Error... %v\n", err)
		return 500, certs
	}

	Log.Printf("CCSAPI 'tlskey' resp StatusCode=%d", resp.StatusCode)
	status = resp.StatusCode
	if resp.StatusCode == 200 {
	
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			Log.Printf("error reading CCSAPI response\n")
			return 500, certs
		}
		//Log.Printf("CCSAPI raw response=%s\n", body)
		err := parse_getCert_Response(body, &certs)
		if err != nil {
			Log.Printf("error parsing ccsapi response\n")
			return 500, certs
		}
		return status, certs  //status == 200
	}
	return status, certs  // status != 200
}


func parse_getCert_Response(body []byte, resp *GetCertResp) error{
	
	err := json.Unmarshal(body, resp)
	if err != nil {
		Log.Println("parse_getCert_Response: error=%v", err)
		return err
	}
	s := fmt.Sprintf("parse_getCert_Response: cert=%s key=%s ", resp.User_cert, resp.User_key)
    	Log.Printf("Retrieved certs from CCSAPI: %v", s)		

	// attempt to create TLS cert based on user's cert and key	
	// TODO this is repeated later, so probably need to be removed from here 	
	_, er := tls.X509KeyPair([]byte(resp.User_cert), []byte(resp.User_key))
	if er != nil {
		Log.Printf("Error loading client key pair for TLS certificate: %v\n", er)
	} else {
		Log.Println("TLS cert is valid")
	}
	return nil
}
