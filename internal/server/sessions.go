package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (app *application) allUserSessions(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(userIDKey).(int64)
	userSessions, err := app.sessionService.GetUserSessions(uid)
	if err != nil {
		log.Printf("error fetching user "+strconv.Itoa(int(uid))+" sessions: %v", err)
		writeError(w, http.StatusInternalServerError, "there was an error retrieving user sessions. try again later")
		return
	}

	w.Header().Add("Content-Type", "application/json")

	//nolint:errcheck
	json.NewEncoder(w).Encode(userSessions)
}

func (app *application) completedUserSessions(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(userIDKey).(int64)
	completedSessions, err := app.sessionService.GetCompletedSessions(uid)
	if err != nil {
		log.Printf("error fetching user "+strconv.Itoa(int(uid))+" sessions: %v", err)
		writeError(w, http.StatusInternalServerError, "there was an error retrieving user sessions. try again later")
		return
	}

	w.Header().Add("Content-Type", "application/json")

	//nolint:errcheck
	json.NewEncoder(w).Encode(completedSessions)
}
