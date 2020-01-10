package crypt

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestCrypt(t *testing.T) {

	tests := []struct {
		pw, salt, expected string
	}{
		{"9CA39C423D4881A6", "H.", "H.4LpYKngoFKI"},
		{"00", "H.", "H.2jPpg5.obl6"},
		{"9CA39C423D4881A6", "ab", "ab3giAvj0qJkg"},
		{"9CA39C423D4881A6", string([]byte{0x41, 0x00}), "AAn2RRYNo7WcI"},
	}

	for i, tt := range tests {
		pw, _ := hex.DecodeString(tt.pw)
		salt := []byte(tt.salt)

		crypted := Crypt(pw, salt)

		if !bytes.Equal(crypted, []byte(tt.expected)) {
			t.Errorf("case %d: Crypt(%s, %s) = %s, want: %s"+
				"\nact: %v \nexp: %v \n",
				i, tt.pw, tt.salt, string(crypted), tt.expected,
				crypted, []byte(tt.expected))
		}
	}
}
