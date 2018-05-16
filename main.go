package main

import (
	"io"
	"log"
	"net/http"
)

const basePath = "/api"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc(basePath+"/heartbeat", HeartbeatHandler)
	log.Println("Starting Server on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func HeartbeatHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}
