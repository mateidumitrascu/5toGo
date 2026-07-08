// Package sessions implements logic related to the study / work session domain entity
package sessions

import "time"

type Session struct {
	SessionID   int64     `json:"sessionId"`
	UserID      int64     `json:"-"`
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Duration    int       `json:"duration"`
	LocalDate   string    `json:"localDate"`
}

func NewSession(sid int64, uid int64, start time.Time, completed time.Time, duration int, local string) *Session {
	return &Session{
		SessionID:   sid,
		UserID:      uid,
		StartedAt:   start,
		CompletedAt: completed,
		Duration:    duration,
		LocalDate:   local,
	}
}
