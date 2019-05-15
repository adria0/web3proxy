package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"io"
	"net/http"
	"net/url"
	"os"
)

var rpcurl string

func web3proxy(w http.ResponseWriter, r *http.Request) {

	fail := func(err error) {
		w.Header().Set("Content-Type", "text/plain")
		log.Error(err)
		w.Write([]byte(err.Error()))
	}

	req, err := http.NewRequest("POST", rpcurl, r.Body)
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
	if err != nil {
		fail(err)
		return
	}
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if _, err := io.Copy(w, resp.Body); err != nil {
		fail(err)
		return
	}
}

func main() {

	if len(os.Args) != 3 {
		log.Info("Usage ", os.Args[0], " externalurl rpcurl")
		os.Exit(-1)
	}

	accessUrl, err := url.Parse(os.Args[1])
	if err != nil || accessUrl.Scheme != "https" {
		log.Error("Malformed URL ", os.Args[1])
		os.Exit(-1)
	}

	rpcurl = os.Args[2]
	http.HandleFunc(accessUrl.Path, web3proxy)

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(accessUrl.Host),
		Cache:      autocert.DirCache("cache-path"),
	}
	server := &http.Server{
		Addr:      ":https",
		TLSConfig: m.TLSConfig(),
	}

	log.Fatal(server.ListenAndServeTLS("", ""))
}
