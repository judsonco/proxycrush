package main

import (
	"github.com/franela/goreq"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ResponseData struct {
	Message string
}

func ProxyHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := r.FormValue("url")
	res, err := goreq.Request{Uri: url}.Do()
	if err != nil {
		log.Fatal(err)
	}

	// Set the content type to whatever was returned
	w.Header().Add("Content-Type", res.Header.Get("Content-Type"))
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
	}

	switch strings.ToLower(res.Header.Get("Content-Type")) {
	case "image/jpeg":
		w.Write(jpegcrush(body))
	default:
		w.Write(body)
	}
}

func jpegcrush(body []byte) []byte {
	c := "jpegtran"
	if os.Getenv("PROXYCRUSH_JPEGTRAN") != "" {
		c = os.Getenv("PROXYCRUSH_JPEGTRAN")
	}

	cmd := exec.Command(
		c,
		"-optimize",
		"-progressive",
		"-copy",
		"none",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Print(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print(err)
	}

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if _, err := stdin.Write(body); err != nil {
		log.Print(err)
	}

	if err = stdin.Close(); err != nil {
		log.Print(err)
	}

	o, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Print(err)
	}

	if err = cmd.Wait(); err != nil {
		log.Print(err)
	}

	return o
}
