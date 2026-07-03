package sessions

import "fmt"

type SessionStore interface {
	Create(s *Session) (*Session, error)
	FindUserSessions(uid int64) ([]Session, error)
	FindCompletedSessions(uid int64) ([]Session, error)
}

type SessionService struct {
	sessionStore SessionStore
}

func NewSessionService(store SessionStore) *SessionService {
	return &SessionService{sessionStore: store}
}

func (srv *SessionService) GetUserSessions(uid int64) ([]Session, error) {
	s, err := srv.sessionStore.FindUserSessions(uid)
	if err != nil {
		return nil, fmt.Errorf("service error getting user sessions: %w", err)
	}
	return s, nil
}

func (srv *SessionService) GetCompletedSessions(uid int64) ([]Session, error) {
	s, err := srv.sessionStore.FindCompletedSessions(uid)
	if err != nil {
		return nil, fmt.Errorf("service error finding completed user sessions: %w", err)
	}
	return s, nil
}
