package testutil

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/types/entity/memcache"
)

type BrokenMemcache struct {
}

func NewBrokenMemcache() *BrokenMemcache {
	return &BrokenMemcache{}
}

func (mem *BrokenMemcache) Set(item *memcache.Item) error {
	return fmt.Errorf("[memcache dummy error] Set(%v)", item)
}

func (mem *BrokenMemcache) Get(key string) (*memcache.Item, error) {
	return nil, fmt.Errorf("[memcache dummy error] Get(%v)", key)
}

func (mem *BrokenMemcache) Delete(key string) error {
	return fmt.Errorf("[memcache dummy error] Delete(%v)", key)
}
