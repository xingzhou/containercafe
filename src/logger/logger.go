package logger

import (
	"log"
	"os"
	"fmt"
	"bytes"
	"sync"
)

type Log struct{
	logger 	* log.Logger
	buf 	bytes.Buffer
	mutex	sync.Mutex
}

//create a Log object that logs to both simultaneously:
//1- stdout - in text lines
//2- file - in logstash understood json format
func NewLogger(logstash_filepath string) (lg * Log){
	lg = new (Log)

	//init standard logger
	log.SetFlags(log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	log.SetPrefix("hijack_proxy: ")
	log.SetOutput(& lg.buf)

	//open log file fp
	fname := logstash_filepath
	fp, err := os.Create(fname)
	if err != nil{
		fmt.Println("Could not create log file ",fname, " will log to stderr only!")
		lg.logger = nil  //redundant
		return
	}
	fmt.Println("Set ELK logging output to ", fname)

	//create log.Logger to write to file fp
	lg.logger = log.New(fp, "", 0)

	return
}

func (lg * Log) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	lg.Output(msg)
}

func (lg * Log) Println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	lg.Output(msg)
}

func (lg * Log) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	lg.Output(msg)
}

func (lg * Log) Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	lg.Output(msg)
	os.Exit(1)
}

func (lg * Log) Output(msg string){
	// safeguard critical section for accessing lg.buf
	lg.mutex.Lock()
	defer lg.mutex.Unlock()
	defer lg.buf.Reset()

	// write standard log line to lg.buf
	//log.Print(msg)  // default calldepth is 2 which is not useful here
	err := log.Output(3, msg)   // not available in go 1.4.2, available in go 1.5.1
	if err != nil {
		fmt.Println("Error - could not write to logger buf")
		return
	}

	// write full standard log line to stdout
	fmt.Printf(lg.buf.String())

	if lg.logger == nil {
		return
	}
	json_msg := lg.format(msg)  // parse buf, msg is passed for convenience but it is part of buf
	lg.logger.Print(json_msg)
}

// Transform std logger line to json
// log line includes: prefix, date, time, file, line, msg
// Example -
// httproxy: 2015/10/01 20:33:59.510644 httproxy.go:155: Resp Status: 200 OK
func (lg * Log) format(msg string) (json_msg string){
	//TODO: parse system log message in lg.buf and create json_msg
	json_msg = "{ \"@fields\": {\"msg\" : \"" + msg + "\"}}"

	//TODO: (enhancement) parse original user msg for potential k=v or k:v pairs and extract them as json fields
	return
}



