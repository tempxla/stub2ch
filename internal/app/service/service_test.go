package service

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/configs/app/secretcfg"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"strconv"
	"testing"
	"time"
)

func TestDefaultBoardService(t *testing.T) {
	sv, err := DefaultBoardService()
	if err != nil {
		t.Error(err)
	}
	if sv.env.SaltComputeId() != secretcfg.COMPUTE_ID_SALT {
		t.Error("unexpected COMPUTE_ID_SALT")
	}
}

func TestNewBoardService(t *testing.T) {
	var repo *BoardStore
	var env *SysEnv
	sv := NewBoardService(RepoConf(repo), EnvConf(env))
	if sv.repo != repo || sv.env != env {
		t.Errorf("%v", sv)
	}
}

func TestMakeDat_ok(t *testing.T) {
	// Setup
	now := time.Now()
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{
			ThreadKey:    "123",
			Dat:          "1行目\n2行目",
			LastModified: now,
		},
	})
	env := &SysEnv{}
	sv := NewBoardService(RepoConf(repo), EnvConf(env))

	// Exercise
	dat, lastModified, err := sv.MakeDat("news4test", "123")

	// Verify
	if err != nil {
		t.Errorf("dat err: %v", err)
	}
	if !bytes.Equal(dat, []byte("1行目\n2行目")) {
		t.Errorf("dat content err. actual: %v", dat)
	}
	if !lastModified.Equal(now) {
		t.Errorf("lastModified: %v", lastModified)
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
	sv := NewBoardService(RepoConf(repo), EnvConf(env))

	// Exercise
	_, _, err := sv.MakeDat("news4test", "999")

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
	sv := NewBoardService(RepoConf(repo), EnvConf(env))

	// Exercise
	txt, err := sv.MakeSubjectTxt("news4test")

	// Verify
	if err != nil {
		t.Errorf("subject.txt err: %v", err)
	}
	if !bytes.Equal(txt,
		[]byte("222.dat<>YYY \t (200)\n111.dat<>XXX \t (100)\n333.dat<>ZZZ \t (300)\n")) {
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
	sv := NewBoardService(RepoConf(repo), EnvConf(env))

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
	testutil.CleanDatastoreBy(t, ctx, client)

	// Sets the kind for the new entity.
	kind := "Board"
	// Sets the name/ID for the new entity.
	name := "news4test"
	// Creates a Key instance.
	boardKey := datastore.NameKey(kind, name, nil)

	// Creates a Board instance.
	boardEntity := board.Entity{
		Subjects: []board.Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &boardEntity); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	sv := NewBoardService(RepoConf(
		&BoardStore{
			Context: ctx,
			Client:  client,
		}),
		EnvConf(&SysEnv{}),
	)
	stng := testutil.NewSettingStub()

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	_, err = sv.CreateThread(stng, "news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Verify
	// ----------------------------------
	if err != nil {
		t.Errorf("err is %v", err)
	}
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(board.Entity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 1 {
		t.Fatalf("subject count  %d", len(e.Subjects))
	}
	// Verify Board
	subject := e.Subjects[0]
	expectedSubject := board.Subject{
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
	var datList []*dat.Entity
	if _, err := client.GetAll(ctx, query, &datList); err != nil {
		t.Fatalf("Failed to get dat  %v", err)
	}
	if len(datList) != 1 {
		t.Fatalf("dat count  %d", len(datList))
	}
	// Verify Dat
	if bytes.Equal(datList[0].Bytes,
		[]byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ")) {
		t.Fatalf("content of dat  %v", datList[0].Bytes)
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
	testutil.CleanDatastoreBy(t, ctx, client)

	// Sets the kind for the new entity.
	kind := "Board"
	// Sets the name/ID for the new entity.
	name := "news4test"
	// Creates a Key instance.
	boardKey := datastore.NameKey(kind, name, nil)

	// Creates a Board instance.
	boardEntity := board.Entity{
		Subjects: []board.Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &boardEntity); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	sv := NewBoardService(RepoConf(
		&BoardStore{
			Context: ctx,
			Client:  client,
		}),
		EnvConf(&SysEnv{}),
	)

	stng := testutil.NewSettingStub()

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	if _, err := sv.CreateThread(stng, "news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト"); err != nil {
		t.Fatalf("at first err is %v", err)
	}

	// Create another thread.
	now2, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 22:29:01.123", time.Local)
	if _, err := sv.CreateThread(stng, "news4test",
		"テスタ2", "age2", now2, "XYZ", "これはテストスレ2", "スレ立てテスト2"); err != nil {
		t.Fatalf("at second err is %v", err)
	}

	// Verify
	// ----------------------------------
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(board.Entity)
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
	testutil.CleanDatastoreBy(t, ctx, client)

	// Injection
	sv := NewBoardService(RepoConf(
		&BoardStore{
			Context: ctx,
			Client:  client,
		}),
		EnvConf(&SysEnv{}),
	)

	//
	stng := testutil.NewSettingStub()

	// Exercise
	// ----------------------------------
	// Create new thread.
	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	_, err = sv.CreateThread(stng, "news4test",
		"テスタ", "age", now, "ABC", "これはテストスレ", "スレ立てテスト")

	// Verify
	// ----------------------------------
	// Get Board
	if err != datastore.ErrNoSuchEntity {
		t.Errorf("err is %v", err)
	}
}

func TestWriteDat(t *testing.T) {
	// Setup
	// ----------------------------------
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	testutil.CleanDatastoreBy(t, ctx, client)

	// Sets the kind for the new entity.
	kind := "Board"
	// Sets the name/ID for the new entity.
	name := "news4test"
	// Creates a Key instance.
	boardKey := datastore.NameKey(kind, name, nil)

	// Creates a Board instance.
	boardEntity := board.Entity{
		Subjects: []board.Subject{},
	}

	// Saves the new entity.
	if _, err := client.Put(ctx, boardKey, &boardEntity); err != nil {
		t.Fatalf("Failed to save board: %v", err)
	}

	// Injection
	startedAt, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-24 23:26:02.789", time.Local)
	sv := NewBoardService(RepoConf(
		&BoardStore{
			Context: ctx,
			Client:  client,
		}),
		EnvConf(&SysEnv{
			StartedTime:   startedAt,
			ComputeIdSalt: "1x",
		}),
	)

	//
	stng := testutil.NewSettingStub()

	// Create new thread.
	date1, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	threadKey, err := sv.CreateThread(stng, "news4test",
		"テスタ", "age", date1, "ABC", "これはテストスレ", "スレ立てテスト")

	// Exercise
	// ----------------------------------
	sv.WriteDat(stng, "news4test", threadKey, "名前2", "メール2", "id2", "カキ２")

	// Verify
	// ----------------------------------
	if err != nil {
		t.Errorf("err is %v", err)
	}
	// Get Board
	key := datastore.NameKey("Board", "news4test", nil)
	e := new(board.Entity)
	if err := client.Get(ctx, key, e); err != nil {
		t.Fatalf("Failed to get board  %v", err)
	}
	if len(e.Subjects) != 1 {
		t.Fatalf("subject count  %d", len(e.Subjects))
	}
	// Verify Board
	subject := e.Subjects[0]
	expectedSubject := board.Subject{
		ThreadKey:    threadKey,
		ThreadTitle:  "スレ立てテスト",
		MessageCount: 2,
		LastModified: sv.StartedAt(),
	}
	if subject != expectedSubject {
		t.Errorf("Fail: contents of subject. actual %v", subject)
		t.Fatalf("Fail: contents of subject. expect %v", expectedSubject)
	}
	// Get Dat
	ancestor := datastore.NameKey("Board", "news4test", nil)
	query := datastore.NewQuery("Dat").Ancestor(ancestor)
	var datList []*dat.Entity
	if _, err := client.GetAll(ctx, query, &datList); err != nil {
		t.Fatalf("Failed to get dat  %v", err)
	}
	if len(datList) != 1 {
		t.Fatalf("dat count  %d", len(datList))
	}
	// Verify Dat
	dateStr := fmt.Sprintf("%s(%s) %s",
		sv.StartedAt().Format(dat_date_layout),
		week_days_jp[sv.StartedAt().Weekday()],
		sv.StartedAt().Format(dat_time_layout),
	)
	if bytes.Equal(datList[0].Bytes,
		[]byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ"+
			"\n名前2<>メール2<>"+dateStr+" ID:id2<> カキ２ <>")) {
		t.Fatalf("content of dat  %v", datList[0].Bytes)
	}
	startedAtStr := sv.StartedAt().Format("2006-01-02 15:04:05.000")
	if datList[0].LastModified.Format("2006-01-02 15:04:05.000") != startedAtStr {
		t.Errorf("last modified: %v", datList[0].LastModified)
	}
}

func TestUpdateSubjectsWhenWriteDat_age(t *testing.T) {
	// Setup
	t1 := time.Now().Add(time.Duration(-1) * time.Hour)
	t2 := time.Now().Add(time.Duration(-2) * time.Hour)
	t3 := time.Now().Add(time.Duration(-3) * time.Hour)
	board := &board.Entity{
		Subjects: []board.Subject{
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

	//
	stng := testutil.NewSettingStub()

	// Exercise
	resunum, err := updateSubjectsWhenWriteDat(stng, board, threadKey, mail, now)
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
	board := &board.Entity{
		Subjects: []board.Subject{
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

	//
	stng := testutil.NewSettingStub()

	// Exercise
	resnum, err := updateSubjectsWhenWriteDat(stng, board, threadKey, mail, now)
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
	board := &board.Entity{
		Subjects: []board.Subject{},
	}
	threadKey := "888"
	mail := "sage"
	now := time.Now()

	//
	stng := testutil.NewSettingStub()

	// Exercise
	_, err := updateSubjectsWhenWriteDat(stng, board, threadKey, mail, now)

	// Verify
	if err == nil {
		t.Errorf("error is nil")
	}
}

func TestUpdateSubjectsWhenWriteDat_1001(t *testing.T) {
	// Setup
	t1 := time.Now().Add(time.Duration(-1) * time.Hour)
	board := &board.Entity{
		Subjects: []board.Subject{
			{
				ThreadKey:    "123",
				MessageCount: 999,
				LastModified: t1,
			},
		}}
	threadKey := "123"
	mail := "sage"
	now := time.Now()

	//
	stng := testutil.NewSettingStub()

	// Exercise
	resnum, err := updateSubjectsWhenWriteDat(stng, board, threadKey, mail, now)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify
	if len(board.Subjects) != 1 {
		t.Errorf("board count: %v", len(board.Subjects))
	}
	if resnum != 1000 {
		t.Errorf("wrong resnum: %v", resnum)
	}
	if board.Subjects[0].ThreadKey != "123" ||
		board.Subjects[0].MessageCount != 1001 ||
		board.Subjects[0].LastModified != now {
		t.Errorf("board content: %v", board.Subjects)
	}
}

func TestCreateDat(t *testing.T) {
	// Exercise
	date, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)
	dat := createDat("名前", "メール", date, "ABC", "本文", "スレタイ")

	// Verify
	excepted := []byte("名前<>メール<>2019/11/23(土) 22:29:01.123 ID:ABC<> 本文 <>スレタイ\n")
	if !bytes.Equal(dat.Bytes, excepted) {
		t.Fatalf("fail \n actual: %v \n expect: %v", string(dat.Bytes), string(excepted))
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
		"\n名前2<>メール2<>2019/11/24(日) 22:29:01.123 ID:XYZ<> 本文2 <>\n")
	if !bytes.Equal(dat.Bytes, excepted) {
		t.Fatalf("fail \n actual: %v \n expect: %v", string(dat.Bytes), string(excepted))
	}
}

func TestComputeId(t *testing.T) {
	// http://age.s22.xrea.com/talk2ch/id.txt
	now, _ := time.ParseInLocation("2006/01/02", "2019/12/26", time.Local)
	sv := NewBoardService(
		RepoConf(&BoardStore{}),
		EnvConf(&SysEnv{
			StartedTime:   now,
			ComputeIdSalt: "1385643578654298",
		}),
	)

	ipAddr := "110.111.112.113"
	boardName := "newsplus"

	id := sv.ComputeId(ipAddr, boardName)

	// if id != "0hGpPuA0" {
	if id != "WmvlSQ2M" {
		t.Errorf("value: %v", id)
	}
}

func TestStartedAt(t *testing.T) {
	startedAt := time.Now()
	env := &SysEnv{
		StartedTime: startedAt,
	}
	sv := NewBoardService(EnvConf(env))

	if sv.StartedAt() != startedAt {
		t.Errorf("\n1: %v\n2: %v", startedAt, sv.StartedAt())
	}
}
