package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/5fives-to-go/internal/api"
)

// func (app *application) allUserSessions(w http.ResponseWriter, r *http.Request) {
// 	uid := r.Context().Value(userIDKey).(int64)
// 	userSessions, err := app.sessionService.GetUserSessions(uid)
// 	if err != nil {
// 		log.Printf("error fetching user "+strconv.Itoa(int(uid))+" sessions: %v", err)
// 		writeError(w, http.StatusInternalServerError, "there was an error retrieving user sessions. try again later")
// 		return
// 	}
//
// 	w.Header().Add("Content-Type", "application/json")
//
// 	//nolint:errcheck
// 	json.NewEncoder(w).Encode(userSessions)
// }

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

func (app *application) recordCompletedSession(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(userIDKey).(int64)
	var req api.RecordSessionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid complete session request")
		return
	}

	if req.StartedAt.IsZero() || req.CompletedAt.IsZero() || req.Duration == 0 {
		writeError(w, http.StatusBadRequest, "invalid values for complete session request")
		return
	}

	_, err = app.sessionService.RecordSession(uid, &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "there was an error recording the completed session")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *application) recordActiveSession(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(userIDKey).(int64)
	var req api.RecordActiveSessionReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid active session request structure")
		return
	}
	if req.ElapsedSeconds <= 0 {
		writeError(w, http.StatusBadRequest, "invalid active session values")
		return
	}

	_, err = app.sessionService.RecordActiveSession(uid, &req)
	if err != nil {
		log.Printf("%v", err)
		writeError(w, http.StatusInternalServerError, "there was an error processing the request")
		return
	}
	writeMessage(w, http.StatusOK, "active session recorded")
}
