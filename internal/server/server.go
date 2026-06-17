// Package server defines the server's infrastructure and functionality
package server

import (
	"fmt"
	"net/http"
)

var (
	address     = "0.0.0.0"
	port        = "8888"
	fullAddress = address + ":" + port
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Hello from %v\n", r)
	})

	return mux
}

func NewHTTPServer() *http.Server {
	return &http.Server{
		Addr:    fullAddress,
		Handler: NewMux(),
	}
}
