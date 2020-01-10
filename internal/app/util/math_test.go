package util

import (
	"testing"
)

func TestMinInt(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{-1, 1, -1},
		{1, 0, 0},
	}

	for i, tt := range tests {
		actual := MinInt(tt.a, tt.b)
		if actual != tt.expected {
			t.Errorf("case %d: MinInt(%d, %d) = %d, want: %d",
				i, tt.a, tt.b, actual, tt.expected)
		}
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{-1, 1, 1},
		{1, 0, 1},
	}

	for i, tt := range tests {
		actual := MaxInt(tt.a, tt.b)
		if actual != tt.expected {
			t.Errorf("case %d: MaxInt(%d, %d) = %d, want: %d",
				i, tt.a, tt.b, actual, tt.expected)
		}
	}
}
