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

	"github.com/5fives-to-go/internal/api"
	"github.com/5fives-to-go/internal/auth"
	"github.com/5fives-to-go/internal/sessions"
	"github.com/5fives-to-go/internal/token"
	"github.com/5fives-to-go/internal/users"
)

type AuthService interface {
	RegisterUser(username string, password string) (*users.User, string, error)
	LoginUser(username string, password string) (*users.User, string, error)
	CheckToken(t string) (*token.AuthToken, error)
	LogoutUser(t string) error
}

type SessionService interface {
	// GetUserSessions(uid int64) ([]sessions.Session, error)
	GetCompletedSessions(uid int64) ([]sessions.Session, error)
	RecordSession(uid int64, req *api.RecordSessionRequest) (*sessions.Session, error)
}

type application struct {
	appStats       ApplicationStatus
	authService    AuthService
	sessionService SessionService
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
	var userdata api.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request structure")
		return
	}

	if userdata.Username == "" || userdata.Password == "" {
		writeError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	_, authToken, err := app.authService.RegisterUser(userdata.Username, userdata.Password)

	if errors.Is(err, users.ErrUserExists) {
		writeError(w, http.StatusConflict, "username is taken")
		return
	}

	if err != nil {
		log.Printf("registration error: %v", err)
		writeError(w, http.StatusInternalServerError, "there was an error processing your request")
		return
	}

	if authToken == "" {
		writeError(w, http.StatusInternalServerError, "there was an error processing your request")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	//nolint:errcheck
	json.NewEncoder(w).Encode(api.RegistrationResponse{Token: authToken})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	var userdata api.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&userdata)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid request structure")
		return
	}

	if userdata.Username == "" || userdata.Password == "" {
		writeError(w, http.StatusBadRequest, "invalid request data")
		return
	}

	user, authToken, err := app.authService.LoginUser(userdata.Username, userdata.Password)

	if errors.Is(err, auth.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err != nil {
		log.Printf("error logging user in: %v", err)
		writeError(w, http.StatusInternalServerError, "there was an error processing your request")
		return
	}

	if authToken == "" {
		writeError(w, http.StatusInternalServerError, "there was an error processing your request")
		return
	}

	fmt.Printf("user %s logged in\n", user.Username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	//nolint:errcheck
	json.NewEncoder(w).Encode(api.LoginResponse{Token: authToken})
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	tokenHash := r.Context().Value(tokenHashKey).(string)
	err := app.authService.LogoutUser(tokenHash)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "there was an error logging out, try again later")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// HELPERS
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(api.ErrorResponse{ErrorMessage: message})
}
