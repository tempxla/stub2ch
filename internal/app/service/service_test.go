package service

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/service/repository"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"strconv"
	"testing"
	"time"
)

func TestDefaultBoardService(t *testing.T) {
	sv, err := DefaultBoardService()
	if sv == nil || err != nil {
		t.Errorf("DefaultBoardService() = %v, %v", sv, err)
	}
}

func TestNewBoardService(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)

	repo := repository.NewBoardStore(ctx, client)
	env := &SysEnv{}
	mem := NewAlterMemcache(ctx, client)

	sv := NewBoardService(RepoConf(repo), EnvConf(env), AdminConf(repo, mem))

	if sv.repo != repo {
		t.Errorf("sv.repo = %v", sv.repo)
	}
	if sv.env != env {
		t.Errorf("sv.env = %v", sv.env)
	}
	if sv.Admin.repo != repo {
		t.Errorf("sv.Admin.repo = %v", sv.Admin.repo)
	}
	if sv.Admin.mem != mem {
		t.Errorf("sv.Admin.mem = %v", sv.Admin.mem)
	}
}

func TestMakeDat(t *testing.T) {
	// Setup
	now := testutil.NewTimeJST(t, "2020-01-13 20:54:12.123")
	repo := testutil.NewBoardStub("news4test", []testutil.ThreadStub{
		{
			ThreadKey:    "123",
			Dat:          "1行目\n2行目\n",
			LastModified: now,
		},
	})

	sv := NewBoardService(RepoConf(repo))

	tests := []struct {
		boardName, threadKey string
		dat                  []byte
		lastModified         time.Time
		err                  error
	}{
		{"news4test", "123", []byte("1行目\n2行目\n"), now, nil},
		{"news4test", "999", nil, time.Time{}, datastore.ErrNoSuchEntity},
	}

	// Exercise
	for i, tt := range tests {
		dat, lastModified, err := sv.MakeDat(tt.boardName, tt.threadKey)
		if !bytes.Equal(dat, tt.dat) || !lastModified.Equal(tt.lastModified) || err != tt.err {
			t.Errorf("%d: sv.MakeDat(%s, %s) = (%v, %v, %v), want: (%v, %v, %v)",
				i, tt.boardName, tt.threadKey, dat, lastModified, err,
				tt.dat, tt.lastModified, tt.err)
		}
	}
}

func TestMakeSubjectTxt(t *testing.T) {
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

	sv := NewBoardService(RepoConf(repo))

	tests := []struct {
		threadKey string
		txt       []byte
		err       error
	}{
		{"news4test", []byte("222.dat<>YYY \t (200)\n111.dat<>XXX \t (100)\n333.dat<>ZZZ \t (300)\n"), nil},
		{"xxxxxxxxx", nil, datastore.ErrNoSuchEntity},
	}

	// Exercise
	for i, tt := range tests {
		txt, err := sv.MakeSubjectTxt(tt.threadKey)
		if !bytes.Equal(txt, tt.txt) || err != tt.err {
			t.Errorf("%d: sv.MakeSubjectTxt(%v) = (%v, %v), want: (%v, %v)",
				i, tt.threadKey, txt, err, tt.txt, tt.err)
		}
	}
}

func TestCreateThread(t *testing.T) {

	repo := testutil.InitialBoardStub("news4test")
	sv := NewBoardService(RepoConf(repo))

	stng := testutil.NewSettingStub()

	type test struct {
		boardName, name, mail, id, message, title string
		time                                      time.Time
		err                                       error
	}

	makeTestSequence := func(boardName string, count int, basei int, baset time.Time) []test {
		ret := make([]test, count)
		for i := 0; i < count; i++ {
			ret[i] = test{
				boardName: boardName,
				name:      "名前" + strconv.Itoa(basei),
				mail:      "メール" + strconv.Itoa(basei),
				id:        "ABCDEFGH" + strconv.Itoa(basei),
				message:   "メッセージ" + strconv.Itoa(basei),
				title:     "タイトル1" + strconv.Itoa(basei),
				time:      baset.Add(time.Duration(i) * time.Second),
				err:       nil,
			}
		}
		return ret
	}

	tests := [][]test{
		// OK: 2つスレ立てる
		{
			{"news4test", "名前1", "メール1", "ABCDEFGH01", "メッセージ1", "タイトル1",
				testutil.NewTimeJST(t, "2020-01-18 11:45:56.123"),
				nil,
			},
			{"news4test", "名前2", "メール2", "ABCDEFGH02", "メッセージ2", "タイトル2",
				testutil.NewTimeJST(t, "2020-01-18 11:45:57.123"),
				nil,
			},
		},
		// Error: 同時刻でのスレ立て
		{
			{"news4test", "名前2", "メール2", "ABCDEFGH02", "メッセージ2", "タイトル2",
				testutil.NewTimeJST(t, "2020-01-18 11:45:57.123"),
				fmt.Errorf("thread key is duplicate"),
			},
		},
		// OK: 板単位のスレッド数制限まで立て続ける
		makeTestSequence("news4test", stng.STUB_THREAD_COUNT()-2,
			3, testutil.NewTimeJST(t, "2020-01-18 11:45:58.123")),
		// Error: 板単位のスレッド数制限
		{
			{"news4test", "名前2", "メール2", "ABCDEFGH02", "メッセージ2", "タイトル2",
				testutil.NewTimeJST(t, "2020-01-18 12:45:57.123"),
				fmt.Errorf("%d: これ以上スレ立てできません。。。", stng.STUB_THREAD_COUNT()),
			},
		},
	}

	expected := testutil.InitialBoardStub("news4test")
	for i, tts := range tests {
		for j, tt := range tts {

			threadKey, err := sv.CreateThread(stng, tt.boardName,
				tt.name, tt.mail, tt.time, tt.id, tt.message, tt.title)

			if tt.err != nil {
				if err == nil {
					t.Errorf("(%d,%d) err is nil, want: %v", i, j, tt.err)
				}
				continue
			}

			sbj := createSubject(tt.time, tt.title)
			expectedBoardEntity := expected.BoardMap[tt.boardName]
			appendSubject(expectedBoardEntity, sbj)

			// verify return value.
			if err != nil {
				t.Errorf("(%d,%d) Error: %v\n"+
					"boardName=%v, name=%v, mail=%v, time=%v, id=%v, message=%v, title=%v",
					i, j, err,
					tt.boardName, tt.name, tt.mail, tt.time, tt.id, tt.message, tt.title)
			}
			if threadKey != sbj.ThreadKey {
				t.Errorf("(%d,%d) ThreadKey = %v, want: %v \n"+
					"boardName=%v, name=%v, mail=%v, time=%v, id=%v, message=%v, title=%v",
					i, j, threadKey, sbj.ThreadKey,
					tt.boardName, tt.name, tt.mail, tt.time, tt.id, tt.message, tt.title)
			}

			// verify datastore.
			boardEntity := &board.Entity{}
			boardKey := repo.BoardKey(tt.boardName)
			if err := repo.GetBoard(boardKey, boardEntity); err != nil {
				t.Errorf("(%d,%d) GetBoard: %v", i, j, err)
			}
			datEntity := &dat.Entity{}
			if err := repo.GetDat(repo.DatKey(threadKey, boardKey), datEntity); err != nil {
				t.Errorf("(%d,%d) GetDat: %v", i, j, err)
			}
			if !testutil.EqualBoardEntity(t, boardEntity, expectedBoardEntity) {
				t.Errorf("(%d,%d): unexpected BoardEntity: \nact:%v \nexp:%v", i, j, boardEntity, expectedBoardEntity)
			}
			expectedDatEntity := createDat(tt.name, tt.mail, tt.time, tt.id, tt.message, tt.title)
			if !testutil.EqualDatEntity(t, datEntity, expectedDatEntity) {
				t.Errorf("(%d,%d): unexpected DatEntity: \nact:%v \nexp:%v", i, j, datEntity, expectedDatEntity)
			}
		}
	}
}

func TestCreateThread_EntityLimit(t *testing.T) {

	repo := testutil.InitialBoardStub("news4test")
	env := &SysEnv{StartedTime: testutil.NewTimeJST(t, "2020-01-18 18:16:51.345")}
	sv := NewBoardService(RepoConf(repo), EnvConf(env))

	stng := testutil.NewSettingStub()

	threadKey, err := sv.CreateThread(stng, "news4test", "name1", "mail1", testutil.NewTimeJST(t, "2020-01-18 12:45:57.123"), "ABCDEFGH01", "message1", "title1")
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < stng.STUB_WRITE_ENTITY_LIMIT(); i++ {
		sv.WriteDat(stng, "news4test", threadKey, "name2", "", "ABCDEFGH02", "message2")
	}

	_, err = sv.CreateThread(stng, "news4test", "nameN", "mailN", testutil.NewTimeJST(t, "2020-01-18 12:45:57.123"), "ABCDEFGH0N", "messageN", "titleN")
	if err == nil {
		t.Errorf("err is nil, want: %v", fmt.Errorf("%d: 今日はこれ以上スレ立てできません。。。", stng.STUB_WRITE_ENTITY_LIMIT()))
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
	sv := NewBoardService(
		RepoConf(repository.NewBoardStore(ctx, client)),
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
	sv := NewBoardService(
		RepoConf(repository.NewBoardStore(ctx, client)),
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
	sv := NewBoardService(
		RepoConf(repository.NewBoardStore(ctx, client)),
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
	sv := NewBoardService(
		RepoConf(repository.NewBoardStore(ctx, client)),
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
		RepoConf(&repository.BoardStore{}),
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
