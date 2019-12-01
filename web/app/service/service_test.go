package service

import (
	"../config"
	E "../entity"
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
	if keys, err = client.GetAll(ctx, query, []*E.DatEntity{}); err != nil {
		t.Fatalf("Failed clean dat. get dat  %v", err)
	}
	if err := client.DeleteMulti(ctx, keys); err != nil {
		t.Fatalf("Failed clean dat. delete dat  %v", err)
	}
	// delete all Board
	query = datastore.NewQuery("Board").KeysOnly()
	if keys, err = client.GetAll(ctx, query, []*E.BoardEntity{}); err != nil {
		t.Fatalf("Failed clean board. get dat  %v", err)
	}
	if err := client.DeleteMulti(ctx, keys); err != nil {
		t.Fatalf("Failed clean board. delete dat  %v", err)
	}
}

func TestMakeDat(t *testing.T) {
	// Setup
	repo := &testutil.BoardStub{
		DatMap: map[string]map[string]*E.DatEntity{
			"news4test": map[string]*E.DatEntity{
				"123": &E.DatEntity{
					Dat: []byte("1行目\n2行目"),
				},
			},
		},
	}
	sv := NewBoardService(repo)

	// Exercise
	dat, err := sv.MakeDat("news4test", "123")

	// Verify
	if err != nil {
		t.Errorf("dat err: %v", err)
	}
	if dat != "1行目\n2行目" {
		t.Errorf("dat content err. actual: %v", dat)
	}
}

func TestMakeSubjectTxt(t *testing.T) {

	// Setup
	repo := &testutil.BoardStub{
		BoardMap: map[string]*E.BoardEntity{
			"news4test": &E.BoardEntity{Subjects: []E.Subject{
				E.Subject{
					ThreadKey:    "111",
					ThreadTitle:  "XXX",
					MessageCount: 100,
					LastFloat:    time.Now().Add(time.Duration(2) * time.Hour),
				},
				E.Subject{
					ThreadKey:    "222",
					ThreadTitle:  "YYY",
					MessageCount: 200,
					LastFloat:    time.Now().Add(time.Duration(3) * time.Hour),
				},
				E.Subject{
					ThreadKey:    "333",
					ThreadTitle:  "ZZZ",
					MessageCount: 300,
					LastFloat:    time.Now().Add(time.Duration(1) * time.Hour),
				},
			}},
		},
	}
	sv := NewBoardService(repo)

	// Exercise
	txt, err := sv.MakeSubjectTxt("news4test")

	// Verify
	if err != nil {
		t.Errorf("subject.txt err: %v", err)
	}
	if txt != "222.dat<>YYY \t (200)\n111.dat<>XXX \t (100)\n333.dat<>ZZZ \t (300)" {
		t.Errorf("subject.txt actual: %v", txt)
	}
}

func TestCreateNewThreadAtFirst(t *testing.T) {

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
	board := E.BoardEntity{
		Subjects: []E.Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &board); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	sv := NewBoardService(&BoardStore{
		Context: ctx,
		Client:  client,
	})

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	sv.CreateNewThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Verify
	// ----------------------------------
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(E.BoardEntity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 1 {
		t.Fatalf("subject count  %d", len(e.Subjects))
	}
	// Verify Board
	subject := e.Subjects[0]
	expectedSubject := E.Subject{
		ThreadKey:    strconv.FormatInt(now.Unix(), 10),
		ThreadTitle:  "スレ立てテスト",
		MessageCount: 1,
		LastFloat:    now,
		LastModified: now,
	}
	if subject != expectedSubject {
		t.Errorf("Fail: contents of subject. actual %v", subject)
		t.Fatalf("Fail: contents of subject. expect %v", expectedSubject)
	}
	// Get Dat
	ancestor := datastore.NameKey("Board", "news4test", nil)
	query := datastore.NewQuery("Dat").Ancestor(ancestor)
	var datList []*E.DatEntity
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

func TestCreateNewThreadMore(t *testing.T) {

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
	board := E.BoardEntity{
		Subjects: []E.Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &board); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	sv := NewBoardService(&BoardStore{
		Context: ctx,
		Client:  client,
	})

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	sv.CreateNewThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Create another thread.
	now2, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 22:29:01.123", time.Local)
	sv.CreateNewThread("news4test",
		"テスタ2", "age2", now2, "XYZ", "これはテストスレ2", "スレ立てテスト2")

	// Verify
	// ----------------------------------
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(E.BoardEntity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 2 {
		t.Fatalf("subject count  %d", len(e.Subjects))
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
