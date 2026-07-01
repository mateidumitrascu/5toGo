package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/5fives-to-go/internal/auth"
	"github.com/5fives-to-go/internal/users"
)

type AuthService interface {
	RegisterUser(username string, password string) (*users.User, error)
	LoginUser(username string, password string) (*users.User, error)
}

type application struct {
	appStats    ApplicationStatus
	authService AuthService
}

type ApplicationStatus struct {
	startTime       time.Time
	activeConnCount atomic.Int64
}

func (app *application) serverHealthHandler(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(app.appStats.startTime)

	responseString := fmt.Sprintf(">> SERVER HEALTH\n\tUPTIME: %s\n\tACTIVE CONNECTIONS: %d\n", uptime.String(), app.appStats.activeConnCount.Load())

	if _, err := fmt.Fprintf(w, "%s", responseString); err != nil {
		fmt.Printf("Error responding with uptime: %v", err)
	}
}

func (app *application) connState(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		app.appStats.activeConnCount.Add(1)
	case http.StateClosed:
		app.appStats.activeConnCount.Add(-1)
	}
}

// Users handlers

type credentialsBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var userdata credentialsBody
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if userdata.Username == "" || userdata.Password == "" {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	_, err = app.authService.RegisterUser(userdata.Username, userdata.Password)

	if errors.Is(err, users.ErrUserExists) {
		w.WriteHeader(http.StatusConflict)

		//nolint:errcheck
		w.Write([]byte("username is taken"))
		return
	}

	if err != nil {
		log.Printf("registration error: %v", err)
		http.Error(w, "there was an error processing your request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//nolint:errcheck
	w.Write([]byte("user created"))
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	var userdata credentialsBody
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if userdata.Username == "" || userdata.Password == "" {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	user, err := app.authService.LoginUser(userdata.Username, userdata.Password)
	if errors.Is(err, auth.ErrInvalidCredentials) {
		w.WriteHeader(http.StatusUnauthorized)

		//nolint:errcheck
		w.Write([]byte("invalid credentials"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		//nolint:errcheck
		w.Write([]byte("there was an error processing your request"))
		return
	}

	fmt.Printf("user %s logged in\n", user.Username)
	w.WriteHeader(http.StatusOK)

	//nolint:errcheck
	w.Write([]byte("user logged in"))
}
