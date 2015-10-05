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
	file	* os.File
	mutex	sync.Mutex
}

//create a Log object that logs to both simultaneously:
//1- stdout - in text lines
//2- file - in logstash understood json format
func NewLogger(logstash_filepath string) (lg * Log){

	lg = new (Log)

	//create log.Logger to write to buf
	lg.logger = log.New(& lg.buf, "", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)

	//open log file
	fname := logstash_filepath
	fp, err := os.Create(fname)
	if err != nil{
		fmt.Println("Could not create logstash file ",fname, " will log to stdout only!")
		lg.file = nil
		return
	}
	lg.file = fp
	lg.Print("Set logstash logging output to ", fname)

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

	// write standard package log line to lg.buf
	//lg.logger.Print(msg)  // default calldepth is 2 which is not useful here
	err := lg.logger.Output(3, msg)   // not available in go 1.4.2, available in go 1.5.1
	if err != nil {
		fmt.Println("Error - could not write to logger buf")
		return
	}

	// write full standard log line to stdout
	fmt.Print(lg.buf.String())

	if lg.file == nil {
		return
	}
	json_msg := lg.format(msg)  // parse buf, msg is passed for convenience but it is part of buf
	_, err = lg.file.WriteString(json_msg)
	if err !=nil {
		fmt.Print("Error - could not write to logger file")
	}
}

// Transform std logger line to json
// log line includes: prefix, date, time, file, line, msg
// Example -
// 2015/10/01 20:33:59.510644 httproxy.go:155: Resp Status: 200 OK
func (lg * Log) format(msg string) (json_msg string){
	// parse lg.buf
	sl := bytes.Split(lg.buf.Bytes(), []byte(" "))
	var date = []byte("")
	var time = []byte("")
	var file = []byte("")
	var line = []byte("")
	if len(sl) >= 3 {
		// looks like a formatted log line
		date = sl[0]
		time = sl[1]
		file_and_line := sl[2]
		// msg is equal to the rest of lg.buf

		sl2 := bytes.Split(file_and_line, []byte(":"))
		if len(sl2) >= 2 {
			file  = sl2[0]
			line = sl2[1]
		}
	}

	// Create json object to be written to logstash log file
	//json_msg = "{ \"@fields\": {\"msg\" : \"" + msg + "\"}}"
	json_msg = "{" +
		"\"@fields\": {" +
			"\"asctime\": \"" + string(date)+" "+string(time) + "\", " +
			"\"message\": \"" + msg + "\", " +
			"\"filename\": \"" + string(file) + "\", " +
			"\"lineno\": " + string(line) +
			"}, " +
		"\"@message\": \"" + msg + "\"" +
		"}\n"

	//TODO: (enhancement) parse original user msg for potential k=v or k:v pairs and extract them as json fields

	return
}



