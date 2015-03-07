package main

import (
	"github.com/pilu/traffic"
)

var router *traffic.Router

func init() {
	router = traffic.New()
	router.Get("/url", ProxyHandler)
}

func main() {
	router.Run()
}