// Package server defines the server's infrastructure and functionality
package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/5fives-to-go/internal/users"
)

var (
	address     = "0.0.0.0"
	port        = "8888"
	fullAddress = address + ":" + port
)

func NewMux(app *application) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.serverHealthHandler)
	mux.HandleFunc("POST /users", app.registerUser)

	return mux
}

func NewHTTPServer(db *sql.DB) *http.Server {
	userrepo := users.NewUserSQLiteRepo(db)

	app := &application{
		appStats: ApplicationStatus{startTime: time.Now()},
		userRepo: userrepo,
	}

	return &http.Server{
		Addr:         fullAddress,
		Handler:      NewMux(app),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		ConnState:    app.connState,
	}
}
