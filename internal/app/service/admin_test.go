package service

import (
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"github.com/tempxla/stub2ch/internal/app/service/repository"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/memcache"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"io/ioutil"
	"testing"
)

func TestVerifySession_notfound(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)

	admin := &AdminFunction{
		mem: mem,
	}

	mem.Delete(admincfg.LOGIN_COOKIE_NAME)

	// Exercise
	err := admin.VerifySession("x")

	// Verify
	if err == nil {
		t.Errorf(`err is nil. want: cache miss error.`)
	}
}

func TestVerifySession_unmatch(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)
	admin := &AdminFunction{
		mem: mem,
	}

	item := &memcache.Item{
		Key:   admincfg.LOGIN_COOKIE_NAME,
		Value: []byte("XXXX"),
	}

	mem.Set(item)

	// Exercise
	err := admin.VerifySession("x")

	// Verify
	if err == nil {
		t.Errorf(`err is nil. want: unmatch value.`)
	}
}

func TestVerifySession(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)
	admin := &AdminFunction{
		mem: mem,
	}

	item := &memcache.Item{
		Key:   admincfg.LOGIN_COOKIE_NAME,
		Value: []byte("XXXX"),
	}

	mem.Set(item)

	// Exercise
	err := admin.VerifySession("XXXX")

	// Verify
	if err != nil {
		t.Error(err)
	}
}

func TestLogin(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)
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
		t.Error(err)
	}

	if len(sid) < 32 { // 16 byte
		t.Errorf("len(sid) < 32: weakness!! %v", sid)
	}
}

func TestLogin_fail(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)
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
	for i, tt := range tests {
		if _, err := admin.Login(string(tt.pass), string(tt.sig)); err == nil {
			t.Errorf("%d: admin.Login(%s, %s) = nil. want: a error. ", i, tt.pass, tt.sig)
		}
	}
}

func TestLogin_MemcacheError(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := testutil.NewBrokenMemcache()
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

	_, err = admin.Login(string(passphrase), string(base64Sig))
	if err == nil {
		t.Error("err is nil. want: [memcache error] Set()")
	}
}

func TestLogout(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	mem := NewAlterMemcache(ctx, client)
	admin := &AdminFunction{
		mem: mem,
	}

	item := &memcache.Item{
		Key:   admincfg.LOGIN_COOKIE_NAME,
		Value: []byte("xxx"),
	}
	mem.Set(item)

	err := admin.Logout()
	if err != nil {
		t.Error(err)
	}
}

func TestCreateBoard(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	admin := &AdminFunction{
		repo: repository.NewBoardStore(ctx, client),
	}

	err := admin.CreateBoard("news4test")
	if err != nil {
		t.Errorf(`first: admin.CreateBoard("news4test") = %v`, err)
	}

	err = admin.CreateBoard("news4test")
	if err == nil {
		t.Error(`second: admin.CreateBoard("news4test") = nil`)
	}
}

func TestCreateBoard_DatastoreError(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	admin := &AdminFunction{
		repo: testutil.NewBrokenBoardStub(),
	}

	err := admin.CreateBoard("news4test")
	if err == nil {
		t.Error(`admin.CreateBoard("news4test") = nil`)
	}
}

func TestGetWriteCount(t *testing.T) {
	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := testutil.EmptyBoardStub()
	admin := &AdminFunction{
		repo: repo,
	}

	repo.PutBoard(repo.BoardKey("news4test1"), &board.Entity{
		WriteCount: 7,
	})
	repo.PutBoard(repo.BoardKey("news4test2"), &board.Entity{
		WriteCount: 13,
	})

	// Verify
	count, err := admin.GetWriteCount()
	if count != 20 || err != nil {
		t.Errorf("admin.GetWriteCount() = %v, %v. want: 20, nil", count, err)
	}
}

func TestResetWriteCount(t *testing.T) {
	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := testutil.EmptyBoardStub()
	admin := &AdminFunction{
		repo: repo,
	}

	repo.PutBoard(repo.BoardKey("news4test1"), &board.Entity{
		WriteCount: 7,
	})
	repo.PutBoard(repo.BoardKey("news4test2"), &board.Entity{
		WriteCount: 13,
	})

	err := admin.ResetWriteCount()
	if err != nil {
		t.Errorf("admin.ResetWriteCount() = %v", err)
	}

	// Verify
	count, err := admin.GetWriteCount()
	if count != 0 || err != nil {
		t.Errorf("admin.ResetWriteCount() failure ? \n "+
			"admin.GetWriteCount() = %v, %v. want: 20, nil", count, err)
	}

}
