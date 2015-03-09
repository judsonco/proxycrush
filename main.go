package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http/fcgi"
	"runtime"
	"strings"
)

var (
	listen = flag.String("listen", "", "TCP or Unix socket to listen on")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", ProxyHandler)

	flag.Parse()
	var err error

	switch {
	case strings.HasPrefix(*listen, "unix:"):
		listener, err := net.Listen("unix", strings.Replace(*listen, "unix:", "", 1))
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, r)
	case *listen != "":
		listener, err := net.Listen("tcp", *listen)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, r)
	default:
		err = fcgi.Serve(nil, r)
	}

	if err != nil {
		log.Fatal(err)
	}
}
