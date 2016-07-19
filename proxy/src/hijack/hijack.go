package main

import (
    "net/http"
	"strconv"
	"strings"
	"os"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	// _ "net/http/pprof" //for profiling only

	"logger"	// my logger package, should be the first user package to init
	"conf"   	// my conf package
	"handler" 	// my handlers
)

var Log * logger.Log = logger.TeeLog

func init(){
	//logger package TeeLog var initialization and init() will take place before this main package init() is executed
	//conf.LoadEnv() is called in init() of conf package, before this main package init() is executed
	Log.Print(conf.GetVerStr())
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
	http.HandleFunc("/apis", handler.KubeEndpointHandler)
	http.HandleFunc("/apis/", handler.KubeEndpointHandler)
	http.HandleFunc("/version", handler.KubeEndpointHandler)
	http.HandleFunc("/swaggerapi/", handler.KubeEndpointHandler)

	// set prefix patterns for Groups handler
	http.HandleFunc("/groups/", handler.GroupsEndpointHandler)
	http.HandleFunc("/groups", handler.GroupsEndpointHandler)
	
	// kubeinit to initialize kube tenant
	http.HandleFunc("/kubeinit", handler.KubeAdminEndpointHandler)

	//Rely on NGINX to route accepted docker/swarm url paths only to hijackproxy
	http.HandleFunc("/", handler.DockerEndpointHandler)
	Log.Printf("All handlers registered")
	
	
	/*for profiling only
	go func() {
		Log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	*/

	// init server on any interface + listen_port
	var err error
	if conf.IsTlsInbound() {
		Log.Printf("Starting TLS listener service")
		// Here is how started the server earlier, before parsing the user certs:
		//	err = http.ListenAndServeTLS(":"+strconv.Itoa(listen_port), conf.GetServerCertFile(), conf.GetServerKeyFile(), nil)
		
		// Here is the new server setup
		//caCert, err := ioutil.ReadFile(conf.GetServerCertFile())
		caCert, err := ioutil.ReadFile(conf.GetCaCertFile())
		if err != nil {
			Log.Println(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
	
		// Setup HTTPS client
		tlsConfig := &tls.Config{
			ClientCAs: caCertPool,
			// NoClientCert
			// RequestClientCert
			// RequireAnyClientCert
			// VerifyClientCertIfGiven
			// RequireAndVerifyClientCert
			//ClientAuth: tls.NoClientCert,
			//ClientAuth: tls.RequestClientCert,
			//ClientAuth: tls.RequireAnyClientCert,
			// ClientAuth: tls.VerifyClientCertIfGiven,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		tlsConfig.BuildNameToCertificate()
	
		server := &http.Server{
			Addr:      ":"+strconv.Itoa(listen_port),
			TLSConfig: tlsConfig,
		}
	
		server.ListenAndServeTLS(conf.GetServerCertFile(), conf.GetServerKeyFile()) 

	} else {
		Log.Printf("Starting non-TLS listener service")
		err = http.ListenAndServe(":"+strconv.Itoa(listen_port), nil)
	}

	//print something and exit on fatal error
	if err != nil {
		Log.Fatal("Aborting because ListenAndServe could not start: ", err)
	}
}
