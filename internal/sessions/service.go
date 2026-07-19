package sessions

import (
	"fmt"
	"log"
	"time"

	"github.com/5fives-to-go/internal/api"
)

const localDateFormat = "2006-01-02"

type SessionStore interface {
	Create(s *Session) (*Session, error)
	// FindUserSessions(uid int64) ([]Session, error)
	FindByDay(uid int64, day string) ([]Session, error)
	FindActiveSession(uid int64) (*ActiveSession, error)
	FindCompletedSessions(uid int64) ([]Session, error)
	UpdateActiveSession(*ActiveSession) (*ActiveSession, error)
	DeleteActiveSession(uid int64) error
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
	s, err := srv.sessionStore.Create(NewSession(0, uid, req.StartedAt, req.CompletedAt, req.Duration, srv.computeUserToday(uid)))
	if err != nil {
		return nil, fmt.Errorf("service error recording session: %w", err)
	}

	return s, nil
}

func (srv *SessionService) GetDailySessions(uid int64, day time.Time) ([]Session, error) {
	s, err := srv.sessionStore.FindByDay(uid, day.Format(localDateFormat))
	if err != nil {
		return nil, fmt.Errorf("service error finding "+day.String()+" sessions: %w", err)
	}
	return s, nil
}

func (srv *SessionService) GetActiveSession(uid int64) (*ActiveSession, error) {
	s, err := srv.sessionStore.FindActiveSession(uid)
	if err != nil {
		return nil, fmt.Errorf("error fetching active session: %w", err)
	}

	if s == nil {
		return nil, nil
	}

	if s.LocalDate != srv.computeUserToday(uid) {
		err := srv.DropActiveSession(uid)
		if err != nil {
			return nil, err
		}
		log.Printf("session deleted for uid %d\n", uid)
		return nil, nil
	}
	return s, nil
}

func (srv *SessionService) DropActiveSession(uid int64) error {
	err := srv.sessionStore.DeleteActiveSession(uid)
	if err != nil {
		return fmt.Errorf("service error deleting active session: %w", err)
	}
	return nil
}

func (srv *SessionService) RecordActiveSession(uid int64, req *api.RecordActiveSessionReq) (*ActiveSession, error) {
	s, err := srv.sessionStore.UpdateActiveSession(NewActiveSession(0, uid, time.Now().UTC(), req.ElapsedSeconds, time.Now().UTC(), srv.computeUserToday(uid)))
	if err != nil {
		return nil, fmt.Errorf("service error finding completed user sessions: %w", err)
	}
	return s, nil
}

func (srv *SessionService) computeUserToday(uid int64) string {
	return time.Now().In(srv.getUserTimezone(uid)).Format(localDateFormat)
}

func (srv *SessionService) getUserTimezone(uid int64) *time.Location {
	return time.Now().Location()
}
