package main

import (
	"net/http"
)

func main() {
	serveMux := http.ServeMux{}
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Handler: &serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
