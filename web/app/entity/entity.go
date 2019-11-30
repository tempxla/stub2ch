package entity

import (
	"time"
)

// Kind=Board
// Key=BoardName
type BoardEntity struct {
	Subjects Subjects `datastore:",noindex"`
}

type Subjects []Subject

type Subject struct {
	ThreadKey    string
	ThreadTitle  string
	MessageCount int
	LastFloat    time.Time
	LastModified time.Time
}

func (s Subjects) Len() int {
	return len(s)
}

func (s Subjects) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Subjects) Less(i, j int) bool {
	// 降順
	return s[i].LastFloat.After(s[j].LastFloat)
}

// Kind=Dat
// Ancestor=Board
// Key=ThreadKey
type DatEntity struct {
	Dat []byte
}
