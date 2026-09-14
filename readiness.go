package main

import (
	"net/http"
)

func healthz(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("content-type", "  text/plain; charset=utf-8")
	responseWriter.WriteHeader(200)
	bytes := []byte("OK")
	responseWriter.Write(bytes)
}
