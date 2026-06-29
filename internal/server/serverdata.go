package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/5fives-to-go/internal/auth"
	"github.com/5fives-to-go/internal/users"
)

type application struct {
	appStats ApplicationStatus
	userRepo users.UserRepo
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
func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	var userdata struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if userdata.Username == "" || userdata.Password == "" {
		http.Error(w, "invalid request data", http.StatusBadRequest)
		return
	}

	_, err = auth.RegisterUser(userdata.Username, userdata.Password, app.userRepo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user created"))
}
