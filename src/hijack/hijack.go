package main

import (
    "net/http"
	"strconv"
	"strings"
	"os"

	"handler" 	// my handlers
	"conf"   	// my conf package
	"logger"	// my logger package
	"auth"		// my auth package
	"httphelper"// my helper
)

var _LOG_TO_FILE_ = true   //feature flag
var Log * logger.Log

func initLogger(){

	Log = logger.NewLogger( conf.GetLogFilePath() )
	conf.SetLogger(Log)
	handler.SetLogger(Log)
	auth.SetLogger(Log)
	httphelper.SetLogger(Log)

	/* using go standard logger
	log.SetFlags(log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	log.SetPrefix("hijackproxy: ")
	if _LOG_TO_FILE_ {
		fname := conf.GetLogFilePath()
		fp, err := os.Create(fname)
		if err != nil{
			log.Println("Could not create log file ",fname, " will use stderr")
			return
		}
		log.SetOutput(fp)
		log.Println("Set ELK logging output to ", fname)
	}
	*/
}

func main() {
	initLogger()

	Log.Print(conf.GetVerStr())

	conf.LoadEnv()

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
