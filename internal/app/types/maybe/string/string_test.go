package string

import (
	"testing"
)

func TestNothing(t *testing.T) {
	m := Nothing()
	if !m.IsNothing() {
		t.Error("is not nothing")
	}
}

func TestJust(t *testing.T) {
	m := Just("x")
	if !m.IsJust() {
		t.Error("is not just")
	}
}

func TestFromMaybe(t *testing.T) {
	a := Just("x")
	if s := a.FromMaybe("def"); s != "x" {
		t.Errorf("value: %v", s)
	}
	b := Nothing()
	if s := b.FromMaybe("z"); s != "z" {
		t.Errorf("value: %v", s)
	}
}
