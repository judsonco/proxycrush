package main

import (
	"fmt"
	"github.com/franela/goreq"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ResponseData struct {
	Message string
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	url := params["url"]

	res, _ := goreq.Request{Uri: url}.Do()

	fmt.Println("HERERE")

	// Set the content type to whatever was returned
	w.Header().Add("Content-Type", res.Header.Get("Content-Type"))
	body, _ := ioutil.ReadAll(res.Body)

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
		fmt.Println(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	if err = cmd.Start(); err != nil {
		fmt.Println(err)
	}

	if _, err := stdin.Write(body); err != nil {
		fmt.Println(err)
	}

	if err = stdin.Close(); err != nil {
		fmt.Println(err)
	}

	o, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
	}

	if err = cmd.Wait(); err != nil {
		fmt.Println(err)
	}

	return o
}
