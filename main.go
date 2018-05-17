package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const basePath = "/api"

func main() {
	mux := http.NewServeMux()

	port := "80"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}

	mux.HandleFunc(basePath+"/heartbeat", HeartbeatHandler)
	log.Println("Starting Server on Port " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func HeartbeatHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World")
}
