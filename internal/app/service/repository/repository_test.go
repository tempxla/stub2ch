package repository

import (
	"bytes"
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	"github.com/tempxla/stub2ch/tools/app/testutil"
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
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	key := repo.BoardKey("news4test")

	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)

	entity1 := &board.Entity{
		Subjects: []board.Subject{
			{
				ThreadKey:    "0123",
				ThreadTitle:  "xxx",
				MessageCount: 1,
				LastModified: now,
			},
		},
	}
	repo.PutBoard(key, entity1)

	entity2 := &board.Entity{}
	err = repo.GetBoard(key, entity2)

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
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	boardKey := repo.BoardKey("news4test")

	boardEntity := &board.Entity{}
	repo.PutBoard(boardKey, boardEntity)

	datKey := repo.DatKey("012", boardKey)
	datEntity1 := &dat.Entity{
		Bytes: []byte("hogepiyo"),
	}
	repo.PutDat(datKey, datEntity1)

	datEntity2 := &dat.Entity{}
	repo.GetDat(datKey, datEntity2)

	if !bytes.Equal(datEntity1.Bytes, datEntity2.Bytes) {
		t.Errorf("%s vs %s", datEntity1.Bytes, datEntity2.Bytes)
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
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	key := repo.BoardKey("news4test")

	now, _ := time.ParseInLocation("2006-01-02 15:04:05.000",
		"2019-11-23 22:29:01.123", time.Local)

	// Put
	entity1 := &board.Entity{
		Subjects: []board.Subject{
			{
				ThreadKey:    "0123",
				ThreadTitle:  "xxx",
				MessageCount: 1,
				LastModified: now,
			},
		},
	}
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		repo.TxPutBoard(tx, key, entity1)
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	// Get
	entity2 := &board.Entity{}
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		repo.TxPutBoard(tx, key, entity1)
		err = repo.TxGetBoard(tx, key, entity2)
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
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	// Put
	datEntity1 := &dat.Entity{
		Bytes: []byte("hogepiyo"),
	}
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		boardKey := repo.BoardKey("news4test")
		datKey := repo.DatKey("012", boardKey)
		repo.TxPutDat(tx, datKey, datEntity1)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	// Get
	datEntity2 := &dat.Entity{}
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		boardKey := repo.BoardKey("news4test")
		datKey := repo.DatKey("012", boardKey)
		repo.TxGetDat(tx, datKey, datEntity2)
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(datEntity1.Bytes, datEntity2.Bytes) {
		t.Errorf("%s vs %s", datEntity1.Bytes, datEntity2.Bytes)
	}

}
