package repository

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	"github.com/tempxla/stub2ch/tools/app/testutil"
	"testing"
)

// Put したものを Getできるか？
func TestPutAndGetBoard(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	now := testutil.NewTimeJST(t, "2019-11-23 22:29:01.123")
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
	key := repo.BoardKey("news4test")
	// *** Put ***
	if err := repo.PutBoard(key, entity1); err != nil {
		t.Error(err)
	}

	entity2 := &board.Entity{}
	// *** Get ***
	if err := repo.GetBoard(key, entity2); err != nil {
		t.Error(err)
	}

	// Verify
	if !testutil.EqualBoardEntity(t, entity1, entity2) {
		t.Errorf("entity1 != entity2: \n%v \n%v", entity1, entity2)
	}
}

// Put したものを Getできるか？
func TestPutAndGetDat(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	boardKey := repo.BoardKey("news4test")
	boardEntity := &board.Entity{}
	repo.PutBoard(boardKey, boardEntity)

	datKey := repo.DatKey("012", boardKey)
	datEntity1 := &dat.Entity{
		Bytes: []byte("ふがふが"),
	}
	// *** Put ***
	if err := repo.PutDat(datKey, datEntity1); err != nil {
		t.Error(err)
	}

	datEntity2 := &dat.Entity{}
	// *** Get ***
	if err := repo.GetDat(datKey, datEntity2); err != nil {
		t.Error(err)
	}

	// Verify
	if !testutil.EqualDatEntity(t, datEntity1, datEntity2) {
		t.Errorf("datEntity1 = %v, datEntity2 = %v", datEntity1, datEntity2)
	}
}

func TestGetAllBoard(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	entities1 := []*board.Entity{
		{
			Subjects: []board.Subject{
				{
					ThreadKey:    "0123",
					ThreadTitle:  "xxx",
					MessageCount: 1,
					LastModified: testutil.NewTimeJST(t, "2019-11-23 22:29:01.123"),
				},
			},
		},
		{
			Subjects: []board.Subject{
				{
					ThreadKey:    "4567",
					ThreadTitle:  "yyy",
					MessageCount: 999,
					LastModified: testutil.NewTimeJST(t, "2020-01-02 12:34:56.999"),
				},
			},
		},
	}
	for i, entity := range entities1 {
		key := repo.BoardKey(fmt.Sprintf("news4test%d", i))
		if err := repo.PutBoard(key, entity); err != nil {
			t.Error(err)
		}
	}

	entities2 := []*board.Entity{}
	// *** GetAll ***
	if _, err := repo.GetAllBoard(&entities2); err != nil {
		t.Error(err)
	}

	// Verify
	if !testutil.EqualBoardEntitiesAsSet(t,
		func(a *board.Entity, b *board.Entity) bool {
			return a.Subjects[0].ThreadKey == b.Subjects[0].ThreadKey
		},
		entities1, entities2) {

		t.Errorf("entities1 != entities2: \n%v \n%v", entities1, entities2)
	}
}

// Put したものを Getできるか？
func TestTxPutAndGetBoard(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	now := testutil.NewTimeJST(t, "2019-11-23 22:29:01.123")
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
	key := repo.BoardKey("news4test")
	err := repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// *** Put ***
		return repo.TxPutBoard(tx, key, entity1)
	})
	if err != nil {
		t.Error(err)
	}

	entity2 := &board.Entity{}
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// *** Get ***
		return repo.TxGetBoard(tx, key, entity2)
	})
	if err != nil {
		t.Error(err)
	}

	// Verify
	if !testutil.EqualBoardEntity(t, entity1, entity2) {
		t.Errorf("entity1 != entity2: \n%v \n%v", entity1, entity2)
	}
}

// Put したものを Getできるか？
func TestTxPutAndGetDat(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	datEntity1 := &dat.Entity{
		Bytes: []byte("ふがふが"),
	}
	datKey := repo.DatKey("012", repo.BoardKey("news4test"))
	err := repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// *** Put ***
		return repo.TxPutDat(tx, datKey, datEntity1)
	})
	if err != nil {
		t.Error(err)
	}

	datEntity2 := &dat.Entity{}
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		// *** Get ***
		return repo.TxGetDat(tx, datKey, datEntity2)
	})
	if err != nil {
		t.Error(err)
	}

	// Verify
	if !testutil.EqualDatEntity(t, datEntity1, datEntity2) {
		t.Errorf("datEntity1 = %v, datEntity2 = %v", datEntity1, datEntity2)
	}
}

func TestTxPutMultiAndGetAllBoard(t *testing.T) {

	ctx, client := testutil.NewContextAndClient(t)
	testutil.CleanDatastoreBy(t, ctx, client)

	repo := NewBoardStore(ctx, client)

	entities1 := []*board.Entity{
		{
			Subjects: []board.Subject{
				{
					ThreadKey:    "0123",
					ThreadTitle:  "xxx",
					MessageCount: 1,
					LastModified: testutil.NewTimeJST(t, "2019-11-23 22:29:01.123"),
				},
			},
		},
		{
			Subjects: []board.Subject{
				{
					ThreadKey:    "4567",
					ThreadTitle:  "yyy",
					MessageCount: 999,
					LastModified: testutil.NewTimeJST(t, "2020-01-02 12:34:56.999"),
				},
			},
		},
	}
	keys1 := []*board.Key{}
	for i, _ := range entities1 {
		keys1 = append(keys1, repo.BoardKey(fmt.Sprintf("news4test%d", i)))
	}

	// *** PutMulti ***
	err := repo.RunInTransaction(func(tx *datastore.Transaction) error {
		return repo.TxPutMultiBoard(tx, keys1, entities1)
	})
	if err != nil {
		t.Error(err)
	}

	// *** GetAll ***
	entities2 := []*board.Entity{}
	var keys2 []*board.Key
	err = repo.RunInTransaction(func(tx *datastore.Transaction) error {
		keys2, err = repo.TxGetAllBoard(tx, &entities2)
		return err
	})
	if err != nil {
		t.Error(err)
	}

	// *** Verify ***
	// key
	if la, lb := len(keys1), len(keys2); la != lb {
		t.Errorf("len(keys1) = %d, len(keys2) = %d", la, lb)
	}
	if (keys1[0].DSKey.Name == keys2[0].DSKey.Name &&
		keys1[1].DSKey.Name == keys2[1].DSKey.Name) ||
		(keys1[0].DSKey.Name == keys2[1].DSKey.Name &&
			keys1[1].DSKey.Name == keys2[0].DSKey.Name) {
		// ok
	} else {
		t.Errorf("keys1 = %v, keys2 = %v", keys1, keys2)
	}

	// entity
	if !testutil.EqualBoardEntitiesAsSet(t,
		func(a *board.Entity, b *board.Entity) bool {
			return a.Subjects[0].ThreadKey == b.Subjects[0].ThreadKey
		},
		entities1, entities2) {

		t.Errorf("entities1 != entities2: \n%v \n%v", entities1, entities2)
	}
}
