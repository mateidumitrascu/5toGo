// Package auth defines the user authentication logic
package auth

import (
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
