package main

import (
	"log"
	"net/http"
	"ws/server"
)

var (
	addr   = ":2022"
	banner = `websocket chat room , address is %s`
)

func main() {

	log.Printf(banner+"\n", addr)

	server.RegisterHandle()

	log.Fatal(http.ListenAndServe(addr, nil))
}
