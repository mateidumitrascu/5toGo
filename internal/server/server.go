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

	// mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Printf("Hello from %v\n", r)
	// 	_, err := fmt.Fprintln(w, "Your hello request was received")
	// 	if err != nil {
	// 		fmt.Printf("Error responding to client: %+v\n", err)
	// 	}
	// })
	//
	// mux.HandleFunc("GET /shut", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("Received request, waiting")
	// 	time.Sleep(3 * time.Second)
	//
	// 	_, err := fmt.Fprintln(w, "Shutdown test request processed")
	// 	if err != nil {
	// 		fmt.Printf("Error responding to client: %+v\n", err)
	// 	}
	// })

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
