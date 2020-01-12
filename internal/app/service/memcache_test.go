package service

import (
	"bytes"
	"github.com/tempxla/stub2ch/internal/app/types/entity/memcache"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)

	item1 := &memcache.Item{
		Key:        "key1",
		Value:      []byte("ばりゅー"),
		Expiration: time.Duration(30) * time.Minute,
	}

	// *** Set ***
	if err := mem.Set(item1); err != nil {
		t.Error(err)
	}

	// *** Get ***
	item2, err := mem.Get("key1")
	if err != nil {
		t.Error(err)
	}

	// Verify
	if item1.Key != item2.Key {
		t.Errorf("item1.Key = %s, item2.Key = %s", item1.Key, item2.Key)
	}
	if !bytes.Equal(item1.Value, item2.Value) {
		t.Errorf("item1.Value = %v, item2.Value = %v", item1.Value, item2.Value)
	}
	if item1.Expiration != item2.Expiration {
		t.Errorf("item1.Expiration = %v, item2.Expiration = %v", item1.Expiration, item2.Expiration)
	}
}

func TestDelete(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)

	item1 := &memcache.Item{
		Key:        "key1",
		Value:      []byte("ばりゅー"),
		Expiration: time.Duration(30) * time.Minute,
	}

	// *** Set ***
	if err := mem.Set(item1); err != nil {
		t.Error(err)
	}

	// *** Delete ***
	err := mem.Delete("key1")
	if err != nil {
		t.Error(err)
	}

	// Verify
	_, err = mem.Get("key1")
	if err != memcache.ErrCacheMiss {
		t.Error(err)
	}
}

func TestDelete_Nothing(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)

	item1 := &memcache.Item{
		Key:        "key1",
		Value:      []byte("ばりゅー"),
		Expiration: time.Duration(30) * time.Minute,
	}

	// *** Set ***
	if err := mem.Set(item1); err != nil {
		t.Error(err)
	}

	// *** Delete ***
	err := mem.Delete("key1")
	if err != nil {
		t.Error(err)
	}

	// Verify
	err = mem.Delete("key1")
	if err != nil {
		t.Error(err)
	}
}
