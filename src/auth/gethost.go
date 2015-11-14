package auth

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"

	"httphelper"  	// my httphelper
	"conf"  		// my conf package
	"logger"
)

//getHost ccsapi call response
type GetHostResp struct {
	Container_id  	string
	Container_name 	string
	Host 			string
	Swarm			bool    // True if swarm manager is the target
	Mgr_host		string  // swarm manager host:port
	Swarm_tls		bool	// use tls if true in case of swarm, TODO: respect this flag
	Space_id		string  // for Authorization (tenant isolation) in case of swarm
	Namespace		string	// registry namespace for this tenant's org, used for validating images the user can access
	Apikey			string  // apikey cred of caller
}

var Log * logger.Log

func SetLogger(lg * logger.Log){
	Log = lg
}

//forward r header only without body to ccsapi auth endpoint
// return ok=true if resp status is 200, otherwise ok=false
func getHost(r *http.Request, id string) (status int, host GetHostResp){
	new_uri := "http://"+conf.GetCcsapiHost()+conf.GetCcsapiUri()+"getHost/"+id

	if ! AuthHeadersExist(r.Header){
		Log.Println("Auth headers missing. Will NOT invoke CCSAPI to authenticate.")
		status = 500
		return
	}

	req, _ := http.NewRequest("GET", new_uri, nil)
	httphelper.CopyHeader(req.Header, r.Header)  //req.Header = r.Header
	req.URL.Host = conf.GetCcsapiHost()
	//req.Header.Add(conf.GetCcsapiIdHeader(), id)
	//req.Header.Add(conf.GetCcsapiIdTypeHeader(), "Kubernetes")

	client := &http.Client{
		CheckRedirect: nil,
	}
	Log.Println("will invoke CCSAPI to authenticate...")

	resp, err := client.Do(req)
	if (err != nil) {
		Log.Printf("GetHost: Error... %v\n", err)
		return 500, host
	}

	Log.Printf("resp StatusCode=%d", resp.StatusCode)
	status = resp.StatusCode
	if resp.StatusCode == 200 {
		//first check in header
		//node = httphelper.GetHeader(resp.Header, conf.GetCcsapiComputeNodeHeader())
		//second check for json response in body
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			Log.Printf("error reading ccsapi response\n")
			return 500, host
		}
		Log.Printf("ccsapi raw response=%s\n", body)
		err := parse_getHost_Response(body, &host)
		if err != nil {
			Log.Printf("error parsing ccsapi response\n")
			return 500, host
		}
		return status, host   //status == 200
	}
	return status, host  // status != 200
}

func parse_getHost_Response(body []byte, resp *GetHostResp) error{
	//var resp GetHostResp
	err := json.Unmarshal(body, resp)
	if err != nil {
		Log.Println("parse_getHost_Response: error=%v", err)
		return err
	}
	s := fmt.Sprintf("parse_getHost_Response: host=%s container_id=%s ", resp.Host, resp.Container_id)
	if resp.Swarm {
		s = s + fmt.Sprintf("Mgr_host=%s Space_id=%s Swarm_tls=%t Namespace=%s Apikey=%s",
			resp.Mgr_host, resp.Space_id, resp.Swarm_tls, resp.Namespace, resp.Apikey)
	}
	Log.Printf("%s\n", s)
	return nil
}

func AuthHeadersExist(h http.Header) bool {
	if httphelper.GetHeader (h, "X-Auth-Token") != ""{
		return true
	}
	if httphelper.GetHeader (h, "X-Tls-Client-Dn") != ""{
		return true
	}
	return false
}

