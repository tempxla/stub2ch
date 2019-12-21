package service

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"github.com/tempxla/stub2ch/configs/app/config"
	"io/ioutil"
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

	mem.Delete(admincfg.LOGIN_COOKIE_NAME)

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
		Key:   admincfg.LOGIN_COOKIE_NAME,
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
		Key:   admincfg.LOGIN_COOKIE_NAME,
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

func TestLogin(t *testing.T) {
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

	passphrase, err := ioutil.ReadFile("/tmp/pass_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	sid, err := admin.Login(string(passphrase), string(base64Sig))

	if err != nil {
		t.Errorf("%v", err)
	}

	if len(sid) < 32 { // 16 byte
		t.Errorf("weakness %v", sid)
	}
}

func TestLogin_fail(t *testing.T) {
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

	passphrase, err := ioutil.ReadFile("/tmp/pass_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	base64Sig, err := ioutil.ReadFile("/tmp/sig_stub2ch.txt")
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		pass []byte
		sig  []byte
	}{
		{pass: passphrase, sig: []byte("wrong sig")},
		{pass: []byte("wrong pass"), sig: base64Sig},
	}
	for _, ts := range tests {
		if _, err := admin.Login(string(ts.pass), string(ts.sig)); err == nil {
			t.Error("err is nil")
		}
	}
}

func TestLogout(t *testing.T) {
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
		Key:   admincfg.LOGIN_COOKIE_NAME,
		Value: []byte("xxx"),
	}
	mem.Set(item)

	err = admin.Logout()

	if err != nil {
		t.Error(err)
	}
}
