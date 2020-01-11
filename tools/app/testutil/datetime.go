package testutil

import (
	"testing"
	"time"
)

func NewTimeJST(t *testing.T, datetime string) time.Time {

	t.Helper()

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	now, err := time.ParseInLocation("2006-01-02 15:04:05.000", datetime, jst)
	if err != nil {
		t.Fatal(err)
	}

	return now
}
