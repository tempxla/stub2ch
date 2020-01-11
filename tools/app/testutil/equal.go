package testutil

import (
	"bytes"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
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
		t.Errorf("a.WriteCount = %d, b.WriteCount = %d", a.WriteCount, b.WriteCount)
		ret = false
	}

	return ret
}

func EqualBoardEntitiesAsSet(t *testing.T,
	keyFunc func(a *board.Entity, b *board.Entity) bool,
	as []*board.Entity, bs []*board.Entity) bool {

	t.Helper()

	if la, lb := len(as), len(bs); la != lb {
		t.Errorf("len(as) = %d, len(bs) = %d", la, lb)
		return false
	}

	mark := make([]bool, len(as))

	for i, v := range as {
		for j, w := range bs {
			if keyFunc(v, w) {
				if mark[j] {
					t.Errorf("duplicate")
					return false
				} else if !EqualBoardEntity(t, v, w) {
					t.Errorf("%d: v = %v, %d: w = %v", i, v, j, w)
					return false
				} else {
					mark[j] = true
				}
			}
		}
	}

	for _, m := range mark {
		if !m {
			t.Errorf("keyが一致しない?: %v", mark)
			return false
		}
	}

	return true
}

func EqualDatEntity(t *testing.T, a *dat.Entity, b *dat.Entity) bool {
	t.Helper()

	ret := true

	if !bytes.Equal(a.Bytes, b.Bytes) {
		t.Errorf("\na.Byte = %v, \nb.Bytes = %v", a.Bytes, b.Bytes)
		ret = false
	}

	if !a.LastModified.Equal(b.LastModified) {
		t.Errorf("a.LastModified = %v, b.LastModified = %v", a.LastModified, b.LastModified)
		ret = false
	}

	return ret
}
