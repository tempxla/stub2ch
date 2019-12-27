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
