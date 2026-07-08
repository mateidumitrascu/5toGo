package sessions

import (
	"fmt"
	"time"

	"github.com/5fives-to-go/internal/api"
)

const localDateFormat = "2006-01-02"

type SessionStore interface {
	Create(s *Session) (*Session, error)
	// FindUserSessions(uid int64) ([]Session, error)
	FindCompletedSessions(uid int64) ([]Session, error)
}

type SessionService struct {
	sessionStore SessionStore
}

func NewSessionService(store SessionStore) *SessionService {
	return &SessionService{sessionStore: store}
}

//
// func (srv *SessionService) GetUserSessions(uid int64) ([]Session, error) {
// 	s, err := srv.sessionStore.FindUserSessions(uid)
// 	if err != nil {
// 		return nil, fmt.Errorf("service error getting user sessions: %w", err)
// 	}
// 	return s, nil
// }

func (srv *SessionService) GetCompletedSessions(uid int64) ([]Session, error) {
	s, err := srv.sessionStore.FindCompletedSessions(uid)
	if err != nil {
		return nil, fmt.Errorf("service error finding completed user sessions: %w", err)
	}
	return s, nil
}

func (srv *SessionService) RecordSession(uid int64, req *api.RecordSessionRequest) (*Session, error) {
	s, err := srv.sessionStore.Create(NewSession(0, uid, req.StartedAt, req.CompletedAt, req.Duration, time.Now().In(srv.getUserTimezone(uid)).Format(localDateFormat)))
	if err != nil {
		return nil, fmt.Errorf("service errro recording session: %w", err)
	}

	return s, nil
}

func (srv *SessionService) getUserTimezone(uid int64) *time.Location {
	return time.Now().Location()
}
