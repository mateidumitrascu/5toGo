// Package server Server infrastructure and functionality
package server

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/5fives-to-go/internal/auth"
	"github.com/5fives-to-go/internal/sessions"
	"github.com/5fives-to-go/internal/token"
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
	mux.HandleFunc("POST /api/register", app.registerUser)
	mux.HandleFunc("POST /api/login", app.loginUser)
	mux.Handle("POST /api/logout", app.requireAuth(http.HandlerFunc(app.logoutUser)))
	// mux.Handle("GET /api/sessions/all", app.requireAuth(http.HandlerFunc(app.allUserSessions)))
	mux.Handle("GET /api/sessions/completed", app.requireAuth(http.HandlerFunc(app.completedUserSessions)))
	mux.Handle("POST /api/sessions/completed", app.requireAuth(http.HandlerFunc(app.recordCompletedSession)))

	return mux
}

func NewHTTPServer(db *sql.DB) *http.Server {
	userRepo := users.NewUserSQLiteRepo(db)
	tokenRepo := token.NewTokenSQLiteRepo(db)
	sessionRepo := sessions.NewSessionSQLiteRepo(db)

	authService := auth.NewAuthService(userRepo, tokenRepo)
	sessionService := sessions.NewSessionService(sessionRepo)

	app := &application{
		appStats:       ApplicationStatus{startTime: time.Now()},
		authService:    authService,
		sessionService: sessionService,
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
