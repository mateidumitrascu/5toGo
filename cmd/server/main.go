// Package main runs the application (starts the server)
package main

import (
	"github.com/5fives-to-go/internal/server"
)

func main() {
	srv := server.NewHTTPServer()
	srv.ListenAndServe()
}
