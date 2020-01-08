package util

import (
	"testing"
)

func TestMinInt(t *testing.T) {
	if MinInt(-1, 1) != -1 {
		t.Error("err1")
	}
	if MinInt(1, 0) != 0 {
		t.Error("err2")
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		a        int
		b        int
		expected int
	}{
		{-1, 1, 1},
		{1, 0, 1},
	}

	for i, tt := range tests {
		if a := MaxInt(tt.a, tt.b); a != tt.expected {
			t.Errorf("case %d: %v", i, a)
		}
	}
}
