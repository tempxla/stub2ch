package crypt

import (
	"encoding/hex"
	"testing"
)

func TestCrypt(t *testing.T) {

	pw, _ := hex.DecodeString("9CA39C423D4881A6")
	salt := []byte("H.")

	x := Crypt(pw, salt)

	expected := "H.4LpYKngoFKI"
	actual := string(x)
	if actual != expected {
		t.Errorf("[%v] %v %v", actual, len(actual), len(expected))
	}

}
