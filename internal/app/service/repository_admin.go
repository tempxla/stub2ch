package service

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
)

type AdminBoardRepository interface {
	CreateBoard(boardName string) (err error)
}

type AdminBoardStore struct {
	repo *BoardStore
}

// 空の板を作成する。
// すでに存在する場合エラーを返す。
func (admin *AdminBoardStore) CreateBoard(boardName string) (err error) {

	key := admin.repo.BoardKey(boardName)
	newEntity := &board.Entity{
		Subjects: []board.Subject{},
	}
	entity := &board.Entity{}
	err = admin.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		err := admin.repo.GetBoard(key, entity)
		if err != nil && err != datastore.ErrNoSuchEntity {
			// datastoreのエラー
			return err
		}
		if err == nil {
			// already exists.
			return fmt.Errorf("entity duplicated: %v", boardName)
		}
		// err == datastore.ErrNoSuchEntity
		return admin.repo.PutBoard(key, newEntity)
	})
	return
}
