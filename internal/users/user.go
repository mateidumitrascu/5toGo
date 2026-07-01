// Package users defines interfaces, structs and methods from the domain of the User entity
package users

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID          int64
	Username     string
	PasswordHash []byte
	CreatedAt    time.Time
}

func (u *User) ValidatePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)); err != nil {
		return false
	}
	return true
}

func NewUser(username string, passwordHash []byte) *User {
	return &User{0, username, passwordHash, time.Time{}}
}
