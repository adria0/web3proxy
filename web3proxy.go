package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
)

var rpcurl string
var debug bool
var opid uint64

func web3proxy(w http.ResponseWriter, r *http.Request) {

	var op uint64

	fail := func(err error) {
		w.Header().Set("Content-Type", "text/plain")
		log.Error(err)
		w.Write([]byte(err.Error()))
	}

	var request io.Reader
	if debug {
		op = atomic.AddUint64(&opid, 1)
		requestBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fail(err)
			return
		}
		log.Debug("> [", op, "] '", string(requestBytes), "'")
		request = bytes.NewReader(requestBytes)
	} else {
		request = r.Body
	}

	req, err := http.NewRequest("POST", rpcurl, request)
	if err != nil {
		fail(err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fail(err)
		return
	}

	defer resp.Body.Close()

	var response io.Reader
	if debug {
		responseBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fail(err)
			return
		}
		log.Debug("< [", op, "] '", string(responseBytes), "'")
		response = bytes.NewReader(responseBytes)
	} else {
		response = resp.Body
	}

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := io.Copy(w, response); err != nil {
		fail(err)
		return
	}
}

func main() {

	debug = os.Getenv("DEBUG") == "1"

	if len(os.Args) != 3 {
		log.Info("Usage ", os.Args[0], " externalurl rpcurl")
		os.Exit(-1)
	}

	accessURL, err := url.Parse(os.Args[1])
	if err != nil || accessURL.Scheme != "https" {
		log.Error("Malformed URL ", os.Args[1])
		os.Exit(-1)
	}

	rpcurl = os.Args[2]

	if debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("externalurl=", accessURL)
		log.Debug("rpcurl=", rpcurl)
	}
	http.HandleFunc(accessURL.Path, web3proxy)

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(accessURL.Host),
		Cache:      autocert.DirCache("cache-path"),
	}
	server := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
	}

	log.Info("Started")
	log.Fatal(server.ListenAndServeTLS("", ""))
}
