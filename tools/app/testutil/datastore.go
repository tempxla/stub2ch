package testutil

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	"github.com/tempxla/stub2ch/internal/app/types/entity/memcache"
	"testing"
)

// Setup Utilitiy
func CleanDatastoreBy(t *testing.T, ctx context.Context, client *datastore.Client) {

	t.Helper()

	// 対象のKIND
	kinds := []string{
		board.KIND,
		dat.KIND,
		memcache.KIND,
	}

	for _, kind := range kinds {
		query := datastore.NewQuery(kind).KeysOnly()
		var keys []*datastore.Key
		var err error
		if keys, err = client.GetAll(ctx, query, nil); err != nil {
			t.Errorf("Error: clean %s. <get> : %v", kind, err)
		}
		if err = client.DeleteMulti(ctx, keys); err != nil {
			t.Errorf("Error: clean %s. <delete> : %v", kind, err)
		}
	}
}

func CleanDatastore(t *testing.T) {

	t.Helper()

	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// 対象のKIND
	kinds := []string{
		board.KIND,
		dat.KIND,
		memcache.KIND,
	}

	for _, kind := range kinds {
		query := datastore.NewQuery(kind).KeysOnly()
		var keys []*datastore.Key
		var err error
		if keys, err = client.GetAll(ctx, query, nil); err != nil {
			t.Errorf("Error: clean %s. <get> : %v", kind, err)
		}
		if err = client.DeleteMulti(ctx, keys); err != nil {
			t.Errorf("Error: clean %s. <delete> : %v", kind, err)
		}
	}
}
