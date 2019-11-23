package main

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/http/httptest"
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

func TestCreateNewThread(t *testing.T) {
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

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

	// Create new thread.
	createNewThread("news4test", "スレ立てテスト", "テスタ", "age", "これはテストスレ")

	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(BoardEntity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 1 {
		t.Fatalf("subject count  %d", len(e.Subjects))
	}
}

func TestCreateDat(t *testing.T) {
	date, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	dat := createDat("名前", "メール", date, "ABC", "本文", "スレタイ")
	excepted := []byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ")
	if !bytes.Equal(dat, excepted) {
		t.Fatalf("fail \n actual: %v \n expect: %v", string(dat), string(excepted))
	}
}

func TestAppendDat(t *testing.T) {
	date1, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	date2, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 22:29:01.123", time.Local)

	dat := createDat("名前", "メール", date1, "ABC", "本文", "スレタイ")

	dat = appendDat(dat, "名前2", "メール2", date2, "XYZ", "本文2")

	excepted := []byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ" +
		"\n名前2<>メール2<>2019/11/24(日) 22:29:01.123 ID:XYZ<> 本文2 <>")
	if !bytes.Equal(dat, excepted) {
		t.Fatalf("fail \n actual: %v \n expect: %v", string(dat), string(excepted))
	}
}
