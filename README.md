# web3proxy
[![Go Report Card](https://goreportcard.com/badge/github.com/adriamb/web3proxy)](https://goreportcard.com/report/github.com/adriamb/web3proxy) [![Build Status](https://travis-ci.org/adriamb/web3proxy.svg?branch=master)](https://travis-ci.org/adriamb/web3proxy) 

A minimalistic web3 authenticated proxy, alla infura.

## Usage

  `web3proxy <https_url> <web3_url>`

where 
  
  - `https_url` is the url to publish the web3 endpoint
  - `web3_url` is where the http web3 is located

for instance `web3proxy https://web3.myhost.com/authentication_token http://localhost:8545`

web3proxy will auto generate the TLS server certificate using https://letsencrypt.org

## Build

You can build from the source or just exec with `docker run -p 443:443 adriamb/web3proxy <https_url> <web3_url>`

