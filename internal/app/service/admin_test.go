package service

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/configs/app/config"
	"testing"
)

func TestVerifySession_notfound(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	mem := &AlterMemcache{
		Client:  client,
		Context: ctx,
	}
	admin := &AdminFunction{
		mem: mem,
	}

	mem.Delete(config.ADMIN_COOKIE_NAME)

	// Exercise
	err = admin.VerifySession("x")

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestVerifySession_unmatch(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	mem := &AlterMemcache{
		Client:  client,
		Context: ctx,
	}
	admin := &AdminFunction{
		mem: mem,
	}

	item := &Item{
		Key:   config.ADMIN_COOKIE_NAME,
		Value: []byte("XXXX"),
	}

	mem.Set(item)

	// Exercise
	err = admin.VerifySession("x")

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
}

func TestVerifySession(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	mem := &AlterMemcache{
		Client:  client,
		Context: ctx,
	}
	admin := &AdminFunction{
		mem: mem,
	}

	item := &Item{
		Key:   config.ADMIN_COOKIE_NAME,
		Value: []byte("XXXX"),
	}

	mem.Set(item)

	// Exercise
	err = admin.VerifySession("XXXX")

	// Verify
	if err != nil {
		t.Error(err)
	}
}
