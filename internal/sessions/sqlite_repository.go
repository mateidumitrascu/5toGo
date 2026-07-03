package sessions

import (
	"database/sql"
	"fmt"
)

type SessionSQLiteRepo struct {
	db *sql.DB
}

const SessionsTable = "sessions"

func NewSessionSQLiteRepo(db *sql.DB) *SessionSQLiteRepo {
	return &SessionSQLiteRepo{db: db}
}

func (sr *SessionSQLiteRepo) Create(s *Session) (*Session, error) {
	result, err := sr.db.Exec("INSERT INTO "+SessionsTable+"(uid, started_at, completed_at, duration) VALUES (?, ?, ?, ?)",
		s.UserID, s.StartedAt, s.CompletedAt, s.Duration)
	if err != nil {
		return nil, fmt.Errorf("error executing session insert: %w", err)
	}
	uid, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last inserted session id: %w", err)
	}
	s.UserID = uid
	return s, nil
}

func (sr *SessionSQLiteRepo) FindUserSessions(uid int64) ([]Session, error) {
	rows, err := sr.db.Query("SELECT session_id, uid, started_at, completed_at, duration FROM "+SessionsTable+" WHERE uid=?", uid)
	if err != nil {
		return nil, fmt.Errorf("error getting user sessions: %w", err)
	}

	//nolint:errcheck
	defer rows.Close()

	userSessions, err := parseRows(rows)
	if err != nil {
		return nil, fmt.Errorf("error finding user sessions: %w", err)
	}
	return userSessions, nil
}

func (sr *SessionSQLiteRepo) FindCompletedSessions(uid int64) ([]Session, error) {
	rows, err := sr.db.Query("SELECT session_id, uid, started_at, completed_at, duration FROM "+SessionsTable+" WHERE uid=? AND completed_at IS NOT NULL", uid)
	if err != nil {
		return nil, fmt.Errorf("error getting user sessions: %w", err)
	}

	//nolint:errcheck
	defer rows.Close()

	userSessions, err := parseRows(rows)
	if err != nil {
		return nil, fmt.Errorf("error finding user sessions: %w", err)
	}
	return userSessions, nil
}

func parseRows(rows *sql.Rows) ([]Session, error) {
	s := []Session{}
	for rows.Next() {
		var sesh Session
		err := rows.Scan(
			&sesh.SessionID,
			&sesh.UserID,
			&sesh.StartedAt,
			&sesh.CompletedAt,
			&sesh.Duration,
		)
		if err != nil {
			return nil, fmt.Errorf("error parsing sessions from rows: %w", err)
		}
		s = append(s, sesh)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error when parsing sessions: %w", err)
	}

	return s, nil
}
