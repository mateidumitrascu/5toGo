// Package auth defines the user authentication logic
package auth

import (
	"errors"
	"fmt"
	"time"

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

func (sv *AuthService) RegisterUser(username string, password string) (*users.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("register user: hash password: %w", err)
	}
	user, err := sv.userStore.Create(users.NewUser(username, hash))
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}
	tokenValue := token.GenerateToken()
	_, err = sv.tokenStore.Create(token.NewAuthToken(token.HashToken(tokenValue), user.UID, time.Now().AddDate(0, 0, 10)))
	if err != nil {
		return user, "", nil
	}

	return user, tokenValue, nil
}

func (sv *AuthService) LoginUser(username string, password string) (*users.User, string, error) {
	user, err := sv.userStore.FindByUsername(username)

	if errors.Is(err, users.ErrUserNotFound) {
		return nil, "", fmt.Errorf("login user: %w", ErrInvalidCredentials)
	}

	if err != nil {
		return nil, "", fmt.Errorf("login user: find by username: %w", err)
	}

	if user.ValidatePassword(password) {
		tokenValue := token.GenerateToken()
		t := token.NewAuthToken(token.HashToken(tokenValue), user.UID, time.Now().AddDate(0, 0, 10))
		_, err := sv.tokenStore.Create(t)
		if err != nil {
			return nil, "", fmt.Errorf("error saving token in auth service: %w", err)
		}
		return user, tokenValue, nil
	}

	return nil, "", fmt.Errorf("login user: password check: %w", ErrInvalidCredentials)
}

func (sv *AuthService) CheckToken(t string) (*token.AuthToken, error) {
	authToken, err := sv.tokenStore.FindToken(t)
	if err != nil {
		return nil, fmt.Errorf("error checking token validity: %w", err)
	}

	if authToken == nil {
		return nil, nil
	}

	if !authToken.Expiry.After(time.Now()) {
		return nil, nil
	}

	return authToken, nil
}
