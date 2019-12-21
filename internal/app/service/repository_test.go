package service

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/configs/app/config"
	. "github.com/tempxla/stub2ch/internal/app/types"
	"testing"
	"time"
)

func TestPutAndGetBoard(t *testing.T) {

	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	sv := NewBoardService(RepoConf(&BoardStore{
		Client:  client,
		Context: ctx,
	}))

	key := sv.repo.BoardKey("news4test")

	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)

	entity1 := &BoardEntity{
		Subjects: []Subject{
			{
				ThreadKey:    "0123",
				ThreadTitle:  "xxx",
				MessageCount: 1,
				LastModified: now,
			},
		},
	}
	sv.repo.PutBoard(key, entity1)

	entity2 := &BoardEntity{}
	err = sv.repo.GetBoard(key, entity2)

	if len(entity1.Subjects) != len(entity2.Subjects) {
		t.Errorf("len is not equal %v vs %v", len(entity1.Subjects), len(entity2.Subjects))
	}

	for i, sbj := range entity1.Subjects {
		if sbj != entity2.Subjects[i] {
			t.Errorf("%v vs %v", sbj, entity2.Subjects[i])
		}
	}

}

func TestPutAndGetDat(t *testing.T) {
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	sv := NewBoardService(RepoConf(&BoardStore{
		Client:  client,
		Context: ctx,
	}))

	boardKey := sv.repo.BoardKey("news4test")

	boardEntity := &BoardEntity{}
	sv.repo.PutBoard(boardKey, boardEntity)

	datKey := sv.repo.DatKey("012", boardKey)
	datEntity1 := &DatEntity{
		Dat: []byte("hogepiyo"),
	}
	sv.repo.PutDat(datKey, datEntity1)

	datEntity2 := &DatEntity{}
	sv.repo.GetDat(datKey, datEntity2)

	if !bytes.Equal(datEntity1.Dat, datEntity2.Dat) {
		t.Errorf("%s vs %s", datEntity1.Dat, datEntity2.Dat)
	}
}

func TestTxPutAndGetBoard(t *testing.T) {

	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	sv := NewBoardService(RepoConf(&BoardStore{
		Client:  client,
		Context: ctx,
	}))

	key := sv.repo.BoardKey("news4test")

	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)

	// Put
	entity1 := &BoardEntity{
		Subjects: []Subject{
			{
				ThreadKey:    "0123",
				ThreadTitle:  "xxx",
				MessageCount: 1,
				LastModified: now,
			},
		},
	}
	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		sv.repo.TxPutBoard(tx, key, entity1)
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	// Get
	entity2 := &BoardEntity{}
	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		sv.repo.TxPutBoard(tx, key, entity1)
		err = sv.repo.TxGetBoard(tx, key, entity2)
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	if len(entity1.Subjects) != len(entity2.Subjects) {
		t.Errorf("len is not equal %v vs %v", len(entity1.Subjects), len(entity2.Subjects))
	}
	for i, sbj := range entity1.Subjects {
		if sbj != entity2.Subjects[i] {
			t.Errorf("%v vs %v", sbj, entity2.Subjects[i])
		}
	}
}

func TestTxPutAndGetDat(t *testing.T) {
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, config.PROJECT_ID)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Clean Datastore
	cleanDatastore(t, ctx, client)

	sv := NewBoardService(RepoConf(&BoardStore{
		Client:  client,
		Context: ctx,
	}))

	// Put
	datEntity1 := &DatEntity{
		Dat: []byte("hogepiyo"),
	}
	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		boardKey := sv.repo.BoardKey("news4test")
		datKey := sv.repo.DatKey("012", boardKey)
		sv.repo.TxPutDat(tx, datKey, datEntity1)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	// Get
	datEntity2 := &DatEntity{}
	err = sv.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		boardKey := sv.repo.BoardKey("news4test")
		datKey := sv.repo.DatKey("012", boardKey)
		sv.repo.TxGetDat(tx, datKey, datEntity2)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(datEntity1.Dat, datEntity2.Dat) {
		t.Errorf("%s vs %s", datEntity1.Dat, datEntity2.Dat)
	}

}
