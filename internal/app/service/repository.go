package service

import (
	"cloud.google.com/go/datastore"
	"context"
	. "github.com/tempxla/stub2ch/internal/app/types"
)

type BoardRepository interface {
	GetBoard(key *datastore.Key, entity *BoardEntity) (err error)
	PutBoard(key *datastore.Key, entity *BoardEntity) (err error)
	GetDat(key *datastore.Key, entity *DatEntity) (err error)
	PutDat(key *datastore.Key, entity *DatEntity) (err error)
	RunInTransaction(func(tx *datastore.Transaction) error) (err error)
	TxGetBoard(tx *datastore.Transaction, key *datastore.Key, entity *BoardEntity) (err error)
	TxPutBoard(tx *datastore.Transaction, key *datastore.Key, entity *BoardEntity) (err error)
	TxGetDat(tx *datastore.Transaction, key *datastore.Key, entity *DatEntity) (err error)
	TxPutDat(tx *datastore.Transaction, key *datastore.Key, entity *DatEntity) (err error)
}

type BoardStore struct {
	Context context.Context
	Client  *datastore.Client
}

func (repo *BoardStore) GetBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	err = repo.Client.Get(repo.Context, key, entity)
	return
}

func (repo *BoardStore) PutBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	_, err = repo.Client.Put(repo.Context, key, entity)
	return
}

func (repo *BoardStore) GetDat(key *datastore.Key, entity *DatEntity) (err error) {
	err = repo.Client.Get(repo.Context, key, entity)
	return
}

func (repo *BoardStore) PutDat(key *datastore.Key, entity *DatEntity) (err error) {
	_, err = repo.Client.Put(repo.Context, key, entity)
	return
}

func (repo *BoardStore) RunInTransaction(f func(tx *datastore.Transaction) error) (err error) {
	_, err = repo.Client.RunInTransaction(repo.Context, f)
	return
}

func (repo *BoardStore) TxGetBoard(tx *datastore.Transaction, key *datastore.Key, entity *BoardEntity) (err error) {
	err = tx.Get(key, entity)
	return
}

func (repo *BoardStore) TxPutBoard(tx *datastore.Transaction, key *datastore.Key, entity *BoardEntity) (err error) {
	_, err = tx.Put(key, entity)
	return
}

func (repo *BoardStore) TxGetDat(tx *datastore.Transaction, key *datastore.Key, entity *DatEntity) (err error) {
	err = tx.Get(key, entity)
	return
}

func (repo *BoardStore) TxPutDat(tx *datastore.Transaction, key *datastore.Key, entity *DatEntity) (err error) {
	_, err = tx.Put(key, entity)
	return
}
