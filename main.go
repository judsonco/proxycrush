package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http/fcgi"
	"os"
	"os/signal"
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
	r := httprouter.New()
	r.GET("/", ProxyHandler)
	flag.Parse()
	var err error

	// Handle common process-killing signals so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	switch {
	case strings.HasPrefix(*listen, "unix:"):
		listener, err := net.Listen("unix", strings.Replace(*listen, "unix:", "", 1))
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		go func(c chan os.Signal) {
			sig := <-c
			log.Printf("Caught signal %s: shutting down.", sig)
			listener.Close()
			os.Exit(0)
		}(sigc)

		err = fcgi.Serve(listener, r)
	case *listen != "":
		listener, err := net.Listen("tcp", *listen)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()
		go func(c chan os.Signal) {
			sig := <-c
			log.Printf("Caught signal %s: shutting down.", sig)
			listener.Close()
			os.Exit(0)
		}(sigc)

		err = fcgi.Serve(listener, r)
	default:
		err = fcgi.Serve(nil, r)
	}

	if err != nil {
		log.Fatal(err)
	}
}
