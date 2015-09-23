package auth

import (
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"httphelper"  	// my httphelper
	"conf"  		// my conf package
)

//getHost response msg
type GetHostResp struct {
	Container_id  	string
	Container_name 	string
	Host 			string
	Swarm			bool    // True if swarm manager is the target
	Mgr_host		string  // swarm manager host:port
	Swarm_tls		bool	// use tls if true in case of swarm, TODO: respect this flag
	Space_id		string  // for Authorization (tenant isolation) in case of swarm
}

//forward r header only without body to ccsapi auth endpoint
// return ok=true if resp status is 200, otherwise ok=false
func getHost(r *http.Request, id string) (ok bool, host GetHostResp){
	new_uri := "http://"+conf.GetCcsapiHost()+conf.GetCcsapiUri()+"getHost/"+id

	req, _ := http.NewRequest("GET", new_uri, nil)
	httphelper.CopyHeader(req.Header, r.Header)  //req.Header = r.Header
	req.URL.Host = conf.GetCcsapiHost()
	//req.Header.Add(conf.GetCcsapiIdHeader(), id)
	//req.Header.Add(conf.GetCcsapiIdTypeHeader(), "Kubernetes")

	client := &http.Client{
		CheckRedirect: nil,
	}
	log.Println("will invoke CCSAPI to authenticate...")

	resp, err := client.Do(req)
	if (err != nil) {
		log.Printf("GetHost: Error... %v\n", err)
		return false, host
	}

	log.Printf("resp StatusCode=%d", resp.StatusCode)
	if resp.StatusCode == 200 {
		ok = true
		//first check in header
		//node = httphelper.GetHeader(resp.Header, conf.GetCcsapiComputeNodeHeader())
		//second check for json response in body
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			log.Printf("error reading ccsapi response\n")
			return false, host
		}
		log.Printf("ccsapi raw response=%s\n", body)
		err := parse_getHost_Response(body, &host)
		if err != nil {
			log.Printf("error parsing ccsapi response\n")
			return false, host
		}
		return true, host
	}
	return false, host
}

func parse_getHost_Response(body []byte, resp *GetHostResp) error{
	//var resp GetHostResp
	err := json.Unmarshal(body, resp)
	if err != nil {
		log.Println("parse_getHost_Response: error=%v", err)
		return err
	}
	s := fmt.Sprintf("parse_getHost_Response: host=%s container_id=%s ", resp.Host, resp.Container_id)
	if resp.Swarm {
		s = s + fmt.Sprintf("Mgr_host=%s Space_id=%s Swarm_tls=%t", resp.Mgr_host, resp.Space_id, resp.Swarm_tls)
	}
	log.Printf("%s\n", s)
	return nil
}

