package testutil

import (
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"testing"
)

func EqualBoardEntity(t *testing.T, a *board.Entity, b *board.Entity) bool {
	t.Helper()

	ret := true

	// Subjects
	if la, lb := len(a.Subjects), len(b.Subjects); la != lb {
		t.Errorf("len(a.Subjects) = %d, len(b.Subjects) = %d", la, lb)
		return false
	}
	for i, v := range a.Subjects {
		w := b.Subjects[i]
		if v.ThreadKey != w.ThreadKey {
			t.Errorf("%d: v.ThreadKey = %s, w.ThreadKey = %s", i, v.ThreadKey, w.ThreadKey)
			ret = false
		}
		if v.ThreadTitle != w.ThreadTitle {
			t.Errorf("%d: v.ThreadTitle = %s, w.ThreadTitle = %s", i, v.ThreadTitle, w.ThreadTitle)
			ret = false
		}
		if v.MessageCount != w.MessageCount {
			t.Errorf("%d: v.MessageCount = %d, w.MessageCount = %d", i, v.MessageCount, w.MessageCount)
			ret = false
		}
		if !v.LastModified.Equal(w.LastModified) {
			t.Errorf("%d: v.LastModified = %v, w.LastModified = %v", i, v.LastModified, w.LastModified)
			ret = false
		}
	}

	// WriteCount
	if a.WriteCount != b.WriteCount {
		ret = false
	}

	return ret
}
