package dat

import (
	"cloud.google.com/go/datastore"
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
	Bytes []byte `datastore:",noindex"`
}
