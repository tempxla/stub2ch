package memcache

import (
	"errors"
	"time"
)

const (
	KIND = "MemcacheItem"
)

var (
	ErrCacheMiss = errors.New("memcache: cache miss")
)

type Item struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}

// memcacheの代替品
// Kind=MemcacheItem
// Key=Key
type MemcacheItem struct {
	Value      []byte        `datastore:",noindex"`
	Expiration time.Duration `datastore:",noindex"`
}
