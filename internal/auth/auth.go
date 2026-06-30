// Package auth defines the user authentication logic
package auth

import (
	"errors"
	"fmt"

	"github.com/5fives-to-go/internal/users"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username string, password string, repo users.UserRepo) (*users.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register user: hash password: %w", err)
	}
	user, err := repo.Create(users.NewUser(username, hash))
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return user, nil
}

func LoginUser(username string, password string, repo users.UserRepo) (*users.User, error) {
	user, err := repo.FindByUsername(username)

	if errors.Is(err, users.ErrUserNotFound) {
		return nil, fmt.Errorf("login user: %w", ErrInvalidCredentials)
	}

	if err != nil {
		return nil, fmt.Errorf("login user: find by username: %w", err)
	}

	if user.ValidatePassword(password) {
		return user, nil
	}

	return nil, fmt.Errorf("login user: password check: %w", ErrInvalidCredentials)
}
