package users

import (
	"database/sql"
	"fmt"
	"time"
)

const UsersTable = "users"

type UserSQLiteRepo struct {
	db *sql.DB
}

func NewUserSQLiteRepo(db *sql.DB) *UserSQLiteRepo {
	return &UserSQLiteRepo{db}
}

func (r *UserSQLiteRepo) Create(u *User) (*User, error) {
	result, err := r.db.Exec("INSERT INTO "+UsersTable+" (username, password_hash, created_at) VALUES (?, ?, ?)", u.Username, u.PasswordHash, time.Now())
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("create user %q: %w", u.Username, ErrUserExists)
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create user: get last insert id: %w", err)
	}
	u.UID = insertID
	return u, nil
}

func (r *UserSQLiteRepo) FindByUsername(username string) (*User, error) {
	rows, err := r.db.Query("SELECT uid, username, password_hash, created_at FROM "+UsersTable+" WHERE username=?", username)
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	//nolint:errcheck
	defer rows.Close()

	users, err := parseUsers(rows)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("find user %q: %w", username, ErrUserNotFound)
	}
	return &users[0], nil
}

func parseUsers(rows *sql.Rows) ([]User, error) {
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UID, &u.Username, &u.PasswordHash, &u.CreatedAt); err != nil {
			return []User{}, fmt.Errorf("error scanning user from row: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return users, nil
}
