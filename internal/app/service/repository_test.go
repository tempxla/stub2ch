package service

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/config"
	. "github.com/tempxla/stub2ch/internal/app/types"
	"testing"
	"time"
)

func TestGetBoard(t *testing.T) {

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

	// kuso code
	if fmt.Sprintf("%v", entity1) != "&{[{0123 xxx 1 2019-11-23 22:29:01.123 +0900 JST}]}" {
		t.Errorf("unexpected  %s", fmt.Sprintf("%v", entity1))
	}

	if fmt.Sprintf("%v", entity1) != fmt.Sprintf("%v", entity2) {
		t.Errorf("put or get failed.\n%v\n%v", entity1, entity2)
	}
}

// Already Test at TestGetBoard.
// func TestPutBoard(t *testing.T) {
// 	return
// }

// func TestGetDat(t *testing.T) {
// 	err = repo.Client.Get(repo.Context, key, entity)
// 	return
// }

// func TestPutDat(t *testing.T) {
// 	_, err = repo.Client.Put(repo.Context, key, entity)
// 	return
// }

// func TestRunInTransaction(t *testing.T) {
// 	_, err = repo.Client.RunInTransaction(repo.Context, f)
// 	return
// }

// func TestTxGetBoard(t *testing.T) {
// 	err = tx.Get(key, entity)
// 	return
// }

// func TestTxPutBoard(t *testing.T) {
// 	_, err = tx.Put(key, entity)
// 	return
// }

// func TestTxGetDat(t *testing.T) {
// 	err = tx.Get(key, entity)
// 	return
// }

// func TestTxPutDat(t *testing.T) {
// 	_, err = tx.Put(key, entity)
// 	return
// }
