package main

import (
	"fmt"
	"github.com/franela/goreq"
	"github.com/pilu/traffic"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type ResponseData struct {
	Message string
}

func ProxyHandler(w traffic.ResponseWriter, r *traffic.Request) {
	url := r.Param("url")

	res, _ := goreq.Request{Uri: url}.Do()

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
	// Create a tmp
	infile, _ := ioutil.TempFile(os.TempDir(), "proxycrush_")
	outfile, _ := ioutil.TempFile(os.TempDir(), "proxycrush_")

	infile.Write(body)

	defer os.Remove(infile.Name())
	defer os.Remove(outfile.Name())

	c := "jpegtran"
	if os.Getenv("PROXYCRUSH_JPEGTRAN") != "" {
		c = os.Getenv("PROXYCRUSH_JPEGTRAN")
	}

	o, _ := exec.Command(
		c,
		"-outfile",
		outfile.Name(),
		"-optimize",
		"-progressive",
		"-copy",
		"none",
		infile.Name(),
	).CombinedOutput()

	fmt.Println(string(o))

	b, _ := ioutil.ReadFile(outfile.Name())

	return b
}
