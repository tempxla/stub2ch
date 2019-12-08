package service

import (
	"../config"
	. "../entity"
	"../testutil"
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"strconv"
	"testing"
	"time"
)

// Setup Utilitiy
func cleanDatastore(t *testing.T, ctx context.Context, client *datastore.Client) {

	t.Helper()

	// delete all Dat
	query := datastore.NewQuery("Dat").KeysOnly()
	var keys []*datastore.Key
	var err error
	if keys, err = client.GetAll(ctx, query, []*DatEntity{}); err != nil {
		t.Fatalf("Failed clean dat. get dat  %v", err)
	}
	if err := client.DeleteMulti(ctx, keys); err != nil {
		t.Fatalf("Failed clean dat. delete dat  %v", err)
	}
	// delete all Board
	query = datastore.NewQuery("Board").KeysOnly()
	if keys, err = client.GetAll(ctx, query, []*BoardEntity{}); err != nil {
		t.Fatalf("Failed clean board. get dat  %v", err)
	}
	if err := client.DeleteMulti(ctx, keys); err != nil {
		t.Fatalf("Failed clean board. delete dat  %v", err)
	}
}

func TestNewBoardService(t *testing.T) {
	var repo *BoardStore
	var env *SysEnv
	sv := NewBoardService(repo, env)
	if sv.repo != repo || sv.env != env {
		t.Errorf("%v", sv)
	}
}

func TestMakeDat_ok(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{
			ThreadKey: "123",
			Dat:       "1行目\n2行目",
		},
	})
	env := &SysEnv{}
	sv := NewBoardService(repo, env)

	// Exercise
	dat, err := sv.MakeDat("news4test", "123")

	// Verify
	if err != nil {
		t.Errorf("dat err: %v", err)
	}
	if !bytes.Equal(dat, []byte("1行目\n2行目")) {
		t.Errorf("dat content err. actual: %v", dat)
	}
}

func TestMakeDat_err(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{
			ThreadKey: "123",
			Dat:       "1行目\n2行目",
		},
	})
	env := &SysEnv{}
	sv := NewBoardService(repo, env)

	// Exercise
	_, err := sv.MakeDat("news4test", "999")

	// Verify
	if err == nil {
		t.Errorf("err is nil")
	}
	if err != datastore.ErrNoSuchEntity {
		t.Errorf("err is %v", err)
	}
}

func TestMakeSubjectTxt_ok(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{
			ThreadKey:    "222",
			ThreadTitle:  "YYY",
			MessageCount: 200,
		},
		{
			ThreadKey:    "111",
			ThreadTitle:  "XXX",
			MessageCount: 100,
		},
		{
			ThreadKey:    "333",
			ThreadTitle:  "ZZZ",
			MessageCount: 300,
		},
	})
	env := &SysEnv{}
	sv := NewBoardService(repo, env)

	// Exercise
	txt, err := sv.MakeSubjectTxt("news4test")

	// Verify
	if err != nil {
		t.Errorf("subject.txt err: %v", err)
	}
	if !bytes.Equal(txt,
		[]byte("222.dat<>YYY \t (200)\n111.dat<>XXX \t (100)\n333.dat<>ZZZ \t (300)")) {
		t.Errorf("subject.txt actual: %v", txt)
	}
}

func TestMakeSubjectTxt_err(t *testing.T) {
	// Setup
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{
			ThreadKey:    "222",
			ThreadTitle:  "YYY",
			MessageCount: 200,
		},
		{
			ThreadKey:    "111",
			ThreadTitle:  "XXX",
			MessageCount: 100,
		},
		{
			ThreadKey:    "333",
			ThreadTitle:  "ZZZ",
			MessageCount: 300,
		},
	})
	env := &SysEnv{}
	sv := NewBoardService(repo, env)

	// Exercise
	_, err := sv.MakeSubjectTxt("news4test1")

	// Verify
	if err != datastore.ErrNoSuchEntity {
		t.Errorf("err is %v", err)
	}
}

func TestCreateNewThread_AtFirst(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	// Sets the kind for the new entity.
	kind := "Board"
	// Sets the name/ID for the new entity.
	name := "news4test"
	// Creates a Key instance.
	boardKey := datastore.NameKey(kind, name, nil)

	// Creates a Board instance.
	board := BoardEntity{
		Subjects: []Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &board); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	sv := NewBoardService(
		&BoardStore{
			Context: ctx,
			Client:  client,
		},
		&SysEnv{},
	)

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	_, err = sv.CreateThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Verify
	// ----------------------------------
	if err != nil {
		t.Errorf("err is %v", err)
	}
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(BoardEntity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 1 {
		t.Fatalf("subject count  %d", len(e.Subjects))
	}
	// Verify Board
	subject := e.Subjects[0]
	expectedSubject := Subject{
		ThreadKey:    strconv.FormatInt(now.Unix(), 10),
		ThreadTitle:  "スレ立てテスト",
		MessageCount: 1,
		LastModified: now,
	}
	if subject != expectedSubject {
		t.Errorf("Fail: contents of subject. actual %v", subject)
		t.Fatalf("Fail: contents of subject. expect %v", expectedSubject)
	}
	// Get Dat
	ancestor := datastore.NameKey("Board", "news4test", nil)
	query := datastore.NewQuery("Dat").Ancestor(ancestor)
	var datList []*DatEntity
	if _, err := client.GetAll(ctx, query, &datList); err != nil {
		t.Fatalf("Failed to get dat  %v", err)
	}
	if len(datList) != 1 {
		t.Fatalf("dat count  %d", len(datList))
	}
	// Verify Dat
	if bytes.Equal(datList[0].Dat,
		[]byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ")) {
		t.Fatalf("content of dat  %v", datList[0].Dat)
	}
}

func TestCreateNewThread_More(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	// Sets the kind for the new entity.
	kind := "Board"
	// Sets the name/ID for the new entity.
	name := "news4test"
	// Creates a Key instance.
	boardKey := datastore.NameKey(kind, name, nil)

	// Creates a Board instance.
	board := BoardEntity{
		Subjects: []Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &board); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	sv := NewBoardService(
		&BoardStore{
			Context: ctx,
			Client:  client,
		},
		&SysEnv{},
	)

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	if _, err := sv.CreateThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト"); err != nil {
		t.Fatalf("at first err is %v", err)
	}

	// Create another thread.
	now2, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 22:29:01.123", time.Local)
	if _, err := sv.CreateThread("news4test",
		"テスタ2", "age2", now2, "XYZ", "これはテストスレ2", "スレ立てテスト2"); err != nil {
		t.Fatalf("at second err is %v", err)
	}

	// Verify
	// ----------------------------------
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(BoardEntity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 2 {
		t.Fatalf("subject count  %d", len(e.Subjects))
	}
}

func TestCreateNewThread_NoSuchBoard(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	// Injection
	sv := NewBoardService(
		&BoardStore{
			Context: ctx,
			Client:  client,
		},
		&SysEnv{},
	)

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	_, err = sv.CreateThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Verify
	// ----------------------------------
	// Get Board
	if err != datastore.ErrNoSuchEntity {
		t.Errorf("err is %v", err)
	}
}

func TestUpdateSubjectsWhenWriteDat_age(t *testing.T) {
	// Setup
	t1 := time.Now().Add(time.Duration(-1) * time.Hour)
	t2 := time.Now().Add(time.Duration(-2) * time.Hour)
	t3 := time.Now().Add(time.Duration(-3) * time.Hour)
	board := &BoardEntity{[]Subject{
		{
			ThreadKey:    "123",
			MessageCount: 100,
			LastModified: t1,
		},
		{
			ThreadKey:    "999",
			MessageCount: 200,
			LastModified: t2,
		},
		{
			ThreadKey:    "456",
			MessageCount: 300,
			LastModified: t3,
		},
	}}
	threadKey := "999"
	mail := ""
	now := time.Now()

	// Exercise
	resunum, err := updateSubjectsWhenWriteDat(board, threadKey, mail, now)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify
	if len(board.Subjects) != 3 {
		t.Errorf("board count: %v", len(board.Subjects))
	}
	if resunum != 201 {
		t.Errorf("wrong resnum: %v", resunum)
	}
	if board.Subjects[0].ThreadKey != "999" ||
		board.Subjects[0].MessageCount != 201 ||
		board.Subjects[0].LastModified != now ||
		board.Subjects[1].ThreadKey != "123" ||
		board.Subjects[1].MessageCount != 100 ||
		board.Subjects[1].LastModified != t1 ||
		board.Subjects[2].ThreadKey != "456" ||
		board.Subjects[2].MessageCount != 300 ||
		board.Subjects[2].LastModified != t3 {
		t.Errorf("board content: %v", board.Subjects)
	}
}

func TestUpdateSubjectsWhenWriteDat_sage(t *testing.T) {
	// Setup
	t1 := time.Now().Add(time.Duration(-1) * time.Hour)
	t2 := time.Now().Add(time.Duration(-2) * time.Hour)
	board := &BoardEntity{[]Subject{
		{
			ThreadKey:    "123",
			MessageCount: 100,
			LastModified: t1,
		},
		{
			ThreadKey:    "999",
			MessageCount: 200,
			LastModified: t2,
		},
	}}
	threadKey := "999"
	mail := "sage"
	now := time.Now()

	// Exercise
	resnum, err := updateSubjectsWhenWriteDat(board, threadKey, mail, now)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify
	if len(board.Subjects) != 2 {
		t.Errorf("board count: %v", len(board.Subjects))
	}
	if resnum != 201 {
		t.Errorf("wrong resnum: %v", resnum)
	}
	if board.Subjects[1].ThreadKey != "999" ||
		board.Subjects[1].MessageCount != 201 ||
		board.Subjects[1].LastModified != now ||
		board.Subjects[0].ThreadKey != "123" ||
		board.Subjects[0].MessageCount != 100 ||
		board.Subjects[0].LastModified != t1 {
		t.Errorf("board content: %v", board.Subjects)
	}
}

func TestUpdateSubjectsWhenWriteDat_fail(t *testing.T) {
	// Setup
	board := &BoardEntity{
		[]Subject{},
	}
	threadKey := "888"
	mail := "sage"
	now := time.Now()

	// Exercise
	_, err := updateSubjectsWhenWriteDat(board, threadKey, mail, now)

	// Verify
	if err == nil {
		t.Errorf("error is nil")
	}
}

func TestCreateDat(t *testing.T) {
	// Exercise
	date, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	dat := createDat("名前", "メール", date, "ABC", "本文", "スレタイ")

	// Verify
	excepted := []byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ")
	if !bytes.Equal(dat.Dat, excepted) {
		t.Fatalf("fail \n actual: %v \n expect: %v", string(dat.Dat), string(excepted))
	}
}

func TestAppendDat(t *testing.T) {
	// Setup
	date1, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	dat := createDat("名前", "メール", date1, "ABC", "本文", "スレタイ")

	// Exercise
	date2, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 22:29:01.123", time.Local)
	appendDat(dat, "名前2", "メール2", date2, "XYZ", "本文2")

	// Verify
	excepted := []byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ" +
		"\n名前2<>メール2<>2019/11/24(日) 22:29:01.123 ID:XYZ<> 本文2 <>")
	if !bytes.Equal(dat.Dat, excepted) {
		t.Fatalf("fail \n actual: %v \n expect: %v", string(dat.Dat), string(excepted))
	}
}

func TestComputeId(t *testing.T) {
	// http://age.s22.xrea.com/talk2ch/id.txt
	now, _ := time.ParseInLocation("2006/01/02", "2019/12/26", time.Local)
	sv := NewBoardService(
		&BoardStore{},
		&SysEnv{
			StartedTime:   now,
			ComputeIdSalt: "1385643578654298",
		},
	)

	ipAddr := "110.111.112.113"
	boardName := "newsplus"

	id := sv.ComputeId(ipAddr, boardName)

	// if id != "0hGpPuA0" {
	if id != "WmvlSQ2M" {
		t.Errorf("value: %v", id)
	}
}
