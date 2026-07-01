// Package validation Validation utils
package validation

import (
	"net/mail"
	"unicode/utf8"
)

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func MinChars(value string, n int) bool {
	return utf8.RuneCount([]byte(value)) >= n
}

// Future password strength, for now only length is checked
