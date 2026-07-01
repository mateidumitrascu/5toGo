package token

import (
	"fmt"
	"testing"
)

func TestTokenGeneration(t *testing.T) {
	for i := range 5 {
		fmt.Printf("token %d: %s\n", i, GenerateToken())
	}
}
