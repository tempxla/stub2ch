package crypt

import (
	"encoding/hex"
	"testing"
)

func TestCrypt(t *testing.T) {

	tests := []struct {
		pw       string
		salt     string
		expected string
	}{
		{"9CA39C423D4881A6", "H.", "H.4LpYKngoFKI"},
		{"00", "H.", "H.2jPpg5.obl6"},
	}

	for _, tt := range tests {
		pw, _ := hex.DecodeString(tt.pw)
		salt := []byte(tt.salt)

		x := Crypt(pw, salt)

		expected := tt.expected
		actual := string(x)
		if actual != expected {
			t.Errorf("[%v] %v %v", actual, len(actual), len(expected))
		}
	}
}
