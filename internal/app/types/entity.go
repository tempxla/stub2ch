package types

import (
	"cloud.google.com/go/datastore"
	"time"
)

const (
	KIND_BOARD = "Board"
	KIND_DAT   = "Dat"
)

type BoardKey struct {
	Key *datastore.Key
}

type DatKey struct {
	Key *datastore.Key
}

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
