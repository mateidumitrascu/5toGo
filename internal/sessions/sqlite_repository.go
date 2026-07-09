package sessions

import (
	"database/sql"
	"fmt"
)

type SessionSQLiteRepo struct {
	db *sql.DB
}

const (
	SessionsTable       = "sessions"
	ActiveSessionsTable = "active_sessions"
)

func NewSessionSQLiteRepo(db *sql.DB) *SessionSQLiteRepo {
	return &SessionSQLiteRepo{db: db}
}

func (sr *SessionSQLiteRepo) Create(s *Session) (*Session, error) {
	result, err := sr.db.Exec("INSERT INTO "+SessionsTable+"(uid, started_at, completed_at, duration, local_date) VALUES (?, ?, ?, ?, ?)",
		s.UserID, s.StartedAt, s.CompletedAt, s.Duration, s.LocalDate)
	if err != nil {
		return nil, fmt.Errorf("error executing session insert: %w", err)
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last inserted session id: %w", err)
	}
	s.SessionID = insertID
	return s, nil
}

func (sr *SessionSQLiteRepo) UpdateActiveSession(a *ActiveSession) (*ActiveSession, error) {
	var id int64
	err := sr.db.QueryRow("INSERT INTO "+ActiveSessionsTable+` (uid, started_at, elapsed_seconds, last_updated, local_date) VALUES (?, ?, ?, ?, ?) 
			ON CONFLICT(uid) DO 
			UPDATE SET elapsed_seconds = excluded.elapsed_seconds, last_updated = excluded.last_updated
			RETURNING active_session_id
		`,
		a.UserID, a.StartedAt, a.ElapsedSeconds, a.LastUpdated, a.LocalDate).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error reading upserted active session id: %w", err)
	}
	a.ActiveSessionID = id
	return a, nil
}

// func (sr *SessionSQLiteRepo) FindUserSessions(uid int64) ([]Session, error) {
// 	rows, err := sr.db.Query("SELECT session_id, uid, started_at, completed_at, duration, local_date FROM "+SessionsTable+" WHERE uid=?", uid)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting user sessions: %w", err)
// 	}
//
// 	//nolint:errcheck
// 	defer rows.Close()
//
// 	userSessions, err := parseRows(rows)
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding user sessions: %w", err)
// 	}
// 	return userSessions, nil
// }

func (sr *SessionSQLiteRepo) FindCompletedSessions(uid int64) ([]Session, error) {
	rows, err := sr.db.Query("SELECT session_id, uid, started_at, completed_at, duration, local_date FROM "+SessionsTable+" WHERE uid=? AND completed_at IS NOT NULL", uid)
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
			&sesh.LocalDate,
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
