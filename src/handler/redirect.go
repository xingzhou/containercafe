package handler

import (
	"log"
	"net/http/httputil"
	"net/http"
	"net"
	"bytes"
	"crypto/tls"

	"conf"
)

//Forward req to server
//when tls_override=true tls is disabled for the current request being processed only
//The override is a directive received from ccsapi getHost, for a swarm request when swarm master does not support tls
func redirect(r *http.Request, body []byte, redirect_host string, redirect_resource_id string,
	rewriteURI func(uri string, resource string) string, tls_override bool) (*http.Response, error, *httputil.ClientConn){

	var cc *httputil.ClientConn

	c , err := net.Dial("tcp", redirect_host)
	if err != nil {
		// handle error
		log.Printf("Error connecting to server %s, %v\n", redirect_host, err)
		return nil,err,nil
	}

	if conf.IsTlsOutbound() && !tls_override{
		cert, er := tls.LoadX509KeyPair(conf.GetClientCertFile(), conf.GetClientKeyFile())
		if er != nil {
			log.Printf("Error loading client key pair, %v\n", er)
			return nil,err,nil
		}
		c_tls := tls.Client(c, &tls.Config{InsecureSkipVerify : true, Certificates : []tls.Certificate{cert}})
		cc = httputil.NewClientConn(c_tls, nil)
	}else{
		cc = httputil.NewClientConn(c, nil)
	}

	req, _ := http.NewRequest(r.Method, "http://"+redirect_host+rewriteURI(r.RequestURI, redirect_resource_id),
				bytes.NewReader(body))
	req.Header = r.Header
	//req.Host = redirect_host
	req.URL.Host = redirect_host

	log.Println("will forward request to server...")
	resp, err := cc.Do(req)

	return resp, err, cc
}
