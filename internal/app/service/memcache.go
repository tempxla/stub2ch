package service

import (
	"cloud.google.com/go/datastore"
	"context"
	"time"
)

type BoardMemcache interface {
	Set(item *Item) error
	Get(key string) (*Item, error)
}

// 見せかけのmemcache. 実態はdatastore.
// Kind=MemcacheItem
// Key=Key
type Item struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}

const (
	kind = "MemcacheItem"
)

type AlterMemcache struct {
	Context context.Context
	Client  *datastore.Client
}

func (mem *AlterMemcache) Set(item *Item) error {
	key := datastore.NameKey(kind, item.Key, nil)
	_, err := mem.Client.Put(mem.Context, key, item)
	return err
}

func (mem *AlterMemcache) Get(key string) (*Item, error) {
	dKey := datastore.NameKey(kind, key, nil)
	dst := &Item{}
	err := mem.Client.Get(mem.Context, dKey, dst)
	item := Item{Key: key, Value: dst.Value}
	return &item, err
}

func (mem *AlterMemcache) DeleteAll() error {
	query := datastore.NewQuery(kind).KeysOnly()
	var dst []*Item
	keys, err := mem.Client.GetAll(mem.Context, query, dst)
	if err != nil {
		return err
	}
	return mem.Client.DeleteMulti(mem.Context, keys)
}
