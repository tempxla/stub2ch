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
	Context context.Context
	Client  *datastore.Client
}

func (mem *AlterMemcache) Set(item *memcache.Item) error {
	memItem := &memcache.MemcacheItem{
		Value:      item.Value,
		Expiration: item.Expiration,
	}
	key := datastore.NameKey(memcache.KIND, item.Key, nil)
	_, err := mem.Client.Put(mem.Context, key, memItem)
	return err
}

func (mem *AlterMemcache) Get(key string) (*memcache.Item, error) {
	dskey := datastore.NameKey(memcache.KIND, key, nil)
	memItem := &memcache.MemcacheItem{}
	err := mem.Client.Get(mem.Context, dskey, memItem)
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
	err := mem.Client.Delete(mem.Context, dskey)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return memcache.ErrCacheMiss
		} else {
			return err
		}
	}
	return nil
}
