package token

import "errors"

var ErrDuplicateTokenFound = errors.New("two tokens with the same hash were found")
