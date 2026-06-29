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

type UserRepo interface {
	Create(u *User) (*User, error)
	// Update(u *User) (*User, error)
	// FindById(uid int64) (*User, error)
	// // TODO: check if it's better to actually return the user after deleting
	// Delete(uid int64) error

	// entity specific methods
	FindByUsername(username string) (*User, error)
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
