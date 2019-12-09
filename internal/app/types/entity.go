package types

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
	LastModified time.Time
}

// Kind=Dat
// Ancestor=Board
// Key=ThreadKey
type DatEntity struct {
	Dat []byte
}
