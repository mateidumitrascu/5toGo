// Package server defines the server's infrastructure and functionality
package server

import (
	"net/http"
	"time"
)

var (
	address     = "0.0.0.0"
	port        = "8888"
	fullAddress = address + ":" + port
)

func NewMux(sd *ServerData) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", sd.serverHealthHandler)

	return mux
}

func NewHTTPServer() *http.Server {
	sd := &ServerData{
		startTime: time.Now(),
	}

	return &http.Server{
		Addr:         fullAddress,
		Handler:      NewMux(sd),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ConnState:    sd.connState,
	}
}
