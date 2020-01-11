package repository

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
)

type AdminBoardRepository interface {
	CreateBoard(boardName string) (err error)
	GetWriteCount() (count int, err error)
	ResetWriteCount() (err error)
}

type AdminBoardStore struct {
	repo BoardRepository
}

func NewAdminBoardStore(repo BoardRepository) *AdminBoardStore {
	return &AdminBoardStore{repo: repo}
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

func (admin *AdminBoardStore) GetWriteCount() (_ int, err error) {

	var entities []*board.Entity
	_, err = admin.repo.GetAllBoard(entities)

	if err != nil {
		return
	}

	count := 0
	for _, entity := range entities {
		count += entity.WriteCount
	}

	return count, nil
}

func (admin *AdminBoardStore) ResetWriteCount() (err error) {

	err = admin.repo.RunInTransaction(func(tx *datastore.Transaction) error {

		var entities []*board.Entity
		keys, err := admin.repo.TxGetAllBoard(tx, entities)
		if err != nil {
			return err
		}

		for i, _ := range keys {
			entities[i].WriteCount = 0
		}

		return admin.repo.TxPutMultiBoard(tx, keys, entities)
	})

	return
}
