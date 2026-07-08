package api

import "time"

type RecordSessionRequest struct {
	StartedAt   time.Time `json:"startedAt"`
	CompletedAt time.Time `json:"completedAt"`
	Duration    int       `json:"duration"`
}
