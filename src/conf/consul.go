package conf

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

type ConsulService struct{
	Node		string
	Address		string
	ServiceID	string
    ServiceName	string
	ServiceTags	[]string
	ServicePort	int
}

func GetServiceHost(service string) (host string){
	consul_host := GetConsulIp()+":"+strconv.Itoa(GetConsulPort())
	consul_url := "http://"+consul_host+"/v1/catalog/service/"+service
	Log.Printf("Will invoke Consul... url=%s", consul_url)
	//invoke Consul to get service metadata
	client := &http.Client{}
	req, err := http.NewRequest("GET",consul_url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil{
		Log.Printf("Error connecting to Consul: %v", err)
		return
	}
	defer resp.Body.Close()
	body,_:=ioutil.ReadAll(resp.Body)
	Log.Printf("Consul raw response: %s", string(body))
	if resp.StatusCode != 200 {
		return
	}

	//parse body json
	var consulService []ConsulService
	err = json.Unmarshal(body, &consulService)
	if err != nil {
		Log.Println("Error in parsing Consul response: error=%v", err)
		return
	}

	//pick one of the endpoints
	host = consulService[0].Address + ":" + strconv.Itoa(consulService[0].ServicePort)
	Log.Printf("Consul response: host=%s", host)
	return
}
