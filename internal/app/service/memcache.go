package service

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/internal/app/types/entity/memcache"
)

type BoardMemcache interface {
	Set(item *memcache.Item) error
	Get(key string) (*memcache.Item, error)
	Delete(key string) error
}

type AlterMemcache struct {
	context context.Context
	client  *datastore.Client
}

func NewAlterMemcache(ctx context.Context, client *datastore.Client) *AlterMemcache {
	return &AlterMemcache{
		context: ctx,
		client:  client,
	}
}

func (mem *AlterMemcache) Set(item *memcache.Item) error {
	memItem := &memcache.MemcacheItem{
		Value:      item.Value,
		Expiration: item.Expiration,
	}
	key := datastore.NameKey(memcache.KIND, item.Key, nil)
	_, err := mem.client.Put(mem.context, key, memItem)
	return err
}

func (mem *AlterMemcache) Get(key string) (*memcache.Item, error) {
	dskey := datastore.NameKey(memcache.KIND, key, nil)
	memItem := &memcache.MemcacheItem{}
	err := mem.client.Get(mem.context, dskey, memItem)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, memcache.ErrCacheMiss
		} else {
			return nil, err
		}
	}
	item := &memcache.Item{
		Key:        key,
		Value:      memItem.Value,
		Expiration: memItem.Expiration,
	}
	return item, err
}

func (mem *AlterMemcache) Delete(key string) error {
	dskey := datastore.NameKey(memcache.KIND, key, nil)
	err := mem.client.Delete(mem.context, dskey)
	if err != nil {
		// if no such entities, err is nil.
		return err
	}
	return nil
}
