// Package token implements token logic for user authentication
package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AuthToken struct {
	Hash   string
	UID    int64
	Expiry time.Time
}

const (
	DefaultTokenSize = 16
)

func NewAuthToken(hash string, uid int64, expiry time.Time) *AuthToken {
	return &AuthToken{
		Hash:   hash,
		UID:    uid,
		Expiry: expiry,
	}
}

func GenerateToken() string {
	var tokenSize int
	tokenSizeEnv := os.Getenv("TOKEN_BYTE_SIZE")

	if tokenSizeEnv == "" {
		tokenSize = DefaultTokenSize
	} else {
		n, err := strconv.Atoi(tokenSizeEnv)
		if err != nil {
			fmt.Printf("error parsing token size from .env: %v\nfalling back to default value for token size", err)
			tokenSize = DefaultTokenSize
		} else {
			tokenSize = n
		}
	}

	b := make([]byte, tokenSize)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func HashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
