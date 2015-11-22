package main

import (
    "net/http"
	"strconv"
	"strings"
	"os"

	"logger"	// my logger package, should be the first user package to init
	"conf"   	// my conf package
	"handler" 	// my handlers
)

var Log * logger.Log = logger.TeeLog

func init(){
	Log.Print(conf.GetVerStr())
	conf.LoadEnv()  //Log is used in LoadEnv()
}

func main() {
	//init() is auto called first before main(), and after the init() of all imported packages

	listen_port := conf.GetDefaultListenPort()
	//parse args
	nargs := len(os.Args)
	if nargs > 1 {
		listen_port, _=strconv.Atoi(os.Args[1])
	}
	Log.Printf("Listening on port %d", listen_port)
	if nargs > 2 {
		conf.Default_redirect_host = os.Args[2]
		Log.Printf("Default upstream server: %s", conf.Default_redirect_host)
	}
	if nargs > 3 {
		if strings.ToLower(os.Args[3]) == "true" {
			conf.SetTlsInbound(true)
			conf.SetTlsOutbound(true)
		}
		if strings.ToLower(os.Args[3]) == "false" {
			conf.SetTlsInbound(false)
			conf.SetTlsOutbound(true)
		}
		Log.Printf("tls setup: Inbound=%t, Outbound=%t", conf.IsTlsInbound(), conf.IsTlsOutbound())
	}

	//register handlers for supported url paths, can't register same path twice

	// Healthcheck handler
	http.HandleFunc("/hjproxy/", handler.HealthEndpointHandler)

	// set prefix pattern for Kubernetes handler
	http.HandleFunc("/api/", handler.KubeEndpointHandler)
	http.HandleFunc("/api", handler.KubeEndpointHandler)
	http.HandleFunc("/swaggerapi/", handler.KubeEndpointHandler)

	//Rely on NGINX to route accepted docker url paths only to hijackproxy
	http.HandleFunc("/", handler.DockerEndpointHandler)

	//handler.TestPatt()  //TEST

	// init server on any interface + listen_port
	var err error
	if conf.IsTlsInbound() {
		err = http.ListenAndServeTLS(":"+strconv.Itoa(listen_port), conf.GetServerCertFile(), conf.GetServerKeyFile(), nil)
	} else {
		err = http.ListenAndServe(":"+strconv.Itoa(listen_port), nil)
	}

	//print something and exit on fatal error
	if err != nil {
		Log.Fatal("Aborting because ListenAndServe could not start: ", err)
	}
}
