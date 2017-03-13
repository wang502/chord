package main

import (
	"net/http"

	"log"

	"github.com/gorilla/mux"
	"github.com/wang502/chord"
)

func main() {
	transporter := chord.NewTransporter()
	server := chord.NewServer("TestNode1", chord.DefaultConfig("http://localhost:4000"), transporter)
	router := mux.NewRouter()
	transporter.Install(server, router)
	log.Fatal(http.ListenAndServe(":4000", router))
}