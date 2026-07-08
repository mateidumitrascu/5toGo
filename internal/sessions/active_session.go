package sessions

import "time"

type ActiveSession struct {
	ActiveSessionID int64     `json:"activeId"`
	UserID          int64     `json:"userId"`
	StartedAt       time.Time `json:"startedAt"`
	ElapsedSeconds  int       `json:"elapsedSeconds"`
	LastUpdated     time.Time `json:"lastUpdated"`
	LocalDate       string    `json:"localDate"`
}

func NewActiveSession(activeID int64, userID int64, started time.Time, elapsed int, lastUpdated time.Time, local string) *ActiveSession {
	return &ActiveSession{activeID, userID, started, elapsed, lastUpdated, local}
}
