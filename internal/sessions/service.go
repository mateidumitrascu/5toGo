package sessions

import "fmt"

type SessionStore interface {
	Create(s *Session) (*Session, error)
	FindUserSessions(uid int64) ([]Session, error)
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
		return nil, fmt.Errorf("error getting user sessions: %w", err)
	}
	return s, nil
}
