package main

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHandleBbsCgi_404(t *testing.T) {

	router := httprouter.New()
	router.GET("/:board/bbs.cgi", handleBbsCgi)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test1/bbs.cgi", nil)
	router.ServeHTTP(writer, request)

	if writer.Code != 404 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestHandleBbsCgi_200(t *testing.T) {

	router := httprouter.New()
	router.GET("/:board/bbs.cgi", handleBbsCgi)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/test/bbs.cgi", nil)
	router.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestCreateNewThreadAtFirst(t *testing.T) {

	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
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
	sv := NewBoardService(&BoardStore{
		ctx:    ctx,
		client: client,
	})

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	sv.createNewThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Verify
	// ----------------------------------
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
	var datList []*DatEntity
	if _, err := client.GetAll(ctx, query, &datList); err != nil {
		t.Fatalf("Failed to get dat  %v", err)
	}
	if len(datList) != 1 {
		t.Fatalf("dat count  %d", len(datList))
	}
	// Verify Dat
	if bytes.Equal(datList[0].Dat, []byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ")) {
		t.Fatalf("content of dat  %v", datList[0].Dat)
	}
}

func TestCreateNewThreadMore(t *testing.T) {

	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
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
	sv := NewBoardService(&BoardStore{
		ctx:    ctx,
		client: client,
	})

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	sv.createNewThread("news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Create another thread.
	now2, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 22:29:01.123", time.Local)
	sv.createNewThread("news4test",
		"テスタ2", "age2", now2, "XYZ", "これはテストスレ2", "スレ立てテスト2")

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

// Setup Utilitiy
func cleanDatastore(t *testing.T, ctx context.Context, client *datastore.Client) {
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
