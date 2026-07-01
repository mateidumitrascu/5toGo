// Package auth defines the user authentication logic
package auth

import (
	"errors"
	"fmt"

	"github.com/5fives-to-go/internal/token"
	"github.com/5fives-to-go/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Create(u *users.User) (*users.User, error)
	FindByUsername(username string) (*users.User, error)
}

type TokenStore interface {
	Create(t *token.AuthToken) (*token.AuthToken, error)
	FindToken(value string) (*token.AuthToken, error)
	DeleteByHash(hash string) error
}

type AuthService struct {
	userStore  UserStore
	tokenStore TokenStore
}

func NewAuthService(us UserStore, ts TokenStore) *AuthService {
	return &AuthService{
		userStore:  us,
		tokenStore: ts,
	}
}

func (authsrv *AuthService) RegisterUser(username string, password string) (*users.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register user: hash password: %w", err)
	}
	user, err := authsrv.userStore.Create(users.NewUser(username, hash))
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}
	return user, nil
}

func (authsrv *AuthService) LoginUser(username string, password string) (*users.User, error) {
	user, err := authsrv.userStore.FindByUsername(username)

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
