package types

import (
	"testing"
)

func TestNothing(t *testing.T) {
	a := Nothing()
	if !IsNothing(a) {
		t.Error("is not nothing")
	}
}

func TestJust(t *testing.T) {
	a := Just("x")
	if !IsJust(a) {
		t.Error("is not just")
	}
	if s := FromJust(a); s != "x" {
		t.Errorf("value: %v", s)
	}
}

func TestFromMaybe(t *testing.T) {
	a := Just("x")
	if s := FromMaybe(a, "def"); s != "x" {
		t.Errorf("value: %v", s)
	}
	b := Nothing()
	if s := FromMaybe(b, "z"); s != "z" {
		t.Errorf("value: %v", s)
	}
}
