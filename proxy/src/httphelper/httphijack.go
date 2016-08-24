package httphelper

import(
	"net"
	"net/http"
	"net/http/httputil"
	"bufio"
	"sync"
	"time"

	"github.com/golang/glog"
)

func InitProxyHijack(w http.ResponseWriter, cc *httputil.ClientConn, req_id string, proto string){
	var cli_conn, srv_conn  net.Conn
	var cli_bufrw, srv_bufrw *bufio.ReadWriter
	var srv_bufr *bufio.Reader
	var err error = nil

	//hijack client conn (act as server on this conn)
	hj, ok := w.(http.Hijacker)
	if !ok {
		glog.Errorf("httproxy doesn't support hijacking\n")
		http.Error(w, "httproxy doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	cli_conn, cli_bufrw, err = hj.Hijack()
	if err != nil {
		glog.Errorf("httproxy hijacking error\n")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


	//hijack server conn (act as client on this conn)
	srv_conn, srv_bufr = cc.Hijack()
	srv_bufrw = bufio.NewReadWriter(srv_bufr, bufio.NewWriter(srv_conn))


	//call TCP hijack loop if upgrade proto is tcp
	if proto == "TCP" {
		tcpHijack(cli_conn, cli_bufrw, srv_conn, srv_bufrw, req_id)
		time.Sleep(100*time.Millisecond) //allow time for go routines to shutdown after hijack completion
	}else{
		glog.Errorf("hijack protocol %s not supported\n", proto)
	}
}

//implement tcp hijack loop forwarding raw tcp messages between cli and srv
func tcpHijack (cli_conn net.Conn, cli_bufrw *bufio.ReadWriter, srv_conn net.Conn, srv_bufrw *bufio.ReadWriter,
	req_id string) {
	//start 2 blocking read/forward loops, BUT exit as soon as one of them exits
	var wg sync.WaitGroup
	defer cli_conn.Close()
	defer srv_conn.Close()

	wg.Add(1) //add 1 only not 2, to proceed when one of the two go routines finishes

	prefix := "> (req id: " + req_id + ")"
	go rwloop(cli_bufrw, srv_bufrw, cli_conn, srv_conn, prefix, &wg, false)

	prefix = "< (req id: " + req_id + ")"
	go rwloop(srv_bufrw, cli_bufrw, srv_conn, cli_conn, prefix, &wg, true)

	//wait until one go routines finishes... the server read loop routine
	//If the 2nd go routine calls wg.Done() at its exit (not our case) you need to increment wg.Add() so no panic happens
	//the 2 conn closures happen at exit of tcpHijack, if the client routine is still running it will receive read error and exit
	wg.Wait()
	//wg.Add(1) Do not wait for the client. Also the client routine should not call wg.Done() on exit

	//allow time for other go routine to flush any data in pipe before sockets are closed
	//this sleep should not be needed now since we exit only when server closes its socket (i.e., server read routine exits)
	//time.Sleep(100*time.Millisecond)

	prefix = "(req id: " + req_id + ")"
	glog.Infof("%s Hijack exit and connections close\n", prefix)
}

func rwloop (src_buf, dest_buf *bufio.ReadWriter, src_conn, dest_conn net.Conn,
	print_prefix string, wg *sync.WaitGroup, server_read_loop bool) {

	if server_read_loop {
		// the caller, tcpHijack(), should exit only if the server is done
		defer wg.Done()
	}

	glog.Infof("%s rwloop started\n", print_prefix)
	//s, err := src_buf.ReadString('\n')
	b, err := src_buf.ReadByte()
	for (err == nil) {
		//_, werr := dest_buf.WriteString(s)
		werr := dest_buf.WriteByte(b)
		dest_buf.Flush()
		if werr != nil {
			glog.Errorf("%s Error writing data: %v\n", print_prefix, werr)
			//raise flag for other loop to exit
			//dest_conn.Close()   // reader on this conn will get error and exit his loop as well
			return
		}
		//s, err = src_buf.ReadString('\n')
		b, err = src_buf.ReadByte()
	}
	if err != nil {
		glog.Errorf("%s Error reading data: %v\n", print_prefix, err)
		return
	}
}

