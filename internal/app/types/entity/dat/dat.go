package dat

import (
	"cloud.google.com/go/datastore"
	"time"
)

const (
	KIND = "Dat"
)

type Key struct {
	DSKey *datastore.Key
}

// Kind=Dat
// Ancestor=Board
// Key=ThreadKey
type Entity struct {
	Bytes        []byte    `datastore:",noindex"`
	LastModified time.Time `datastore:",noindex"`
}
