package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/5fives-to-go/internal/token"
)

type ctxKey string

var (
	userIDKey    ctxKey = "userID"
	tokenHashKey ctxKey = "tokenHash"
)

type TokenValidator interface {
	CheckToken(t string) (bool, error)
}

func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			writeError(w, http.StatusUnauthorized, "no authorization provided")
			return
		}

		t := strings.TrimPrefix(authorization, "Bearer ")
		authToken, err := app.authService.CheckToken(t)
		if err != nil {
			log.Printf("error checking token in middleware: %v", err)
			writeError(w, http.StatusInternalServerError, "there was an error checking your authorization")
			return
		}

		if authToken == nil {
			writeError(w, http.StatusUnauthorized, "invalid authorization")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, authToken.UID)
		ctx = context.WithValue(ctx, tokenHashKey, token.HashToken(t))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
