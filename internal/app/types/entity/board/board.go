package board

import (
	"cloud.google.com/go/datastore"
	"time"
)

const (
	KIND = "Board"
)

type Key struct {
	DSKey *datastore.Key
}

// Kind=Board
// Key=BoardName
type Entity struct {
	Subjects   []Subject `datastore:",noindex"`
	WriteCount int       `datastore:",noindex"`
}

type Subject struct {
	ThreadKey    string    `datastore:",noindex"`
	ThreadTitle  string    `datastore:",noindex"`
	MessageCount int       `datastore:",noindex"`
	LastModified time.Time `datastore:",noindex"` // dat落ちとかで使う予定
}
