package server

import (
	"log"
	"net/http"
)

func Init() {
	mux := http.NewServeMux()

	mux.HandleFunc("/account", getAccount)
	log.Fatal(http.ListenAndServe("localhost:3000", mux))
}
