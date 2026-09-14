package main

import (
	"net/http"
)

func main() {
	serveMux := http.ServeMux{}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	serveMux.Handle("/app/", handler)
	serveMux.Handle("/app/assets/logo.png", handler)
	serveMux.HandleFunc("/healthz", healthz)

	server := http.Server{
		Handler: &serveMux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
