package service

import (
	"cloud.google.com/go/datastore"
	"context"
	. "github.com/tempxla/stub2ch/internal/app/types"
)

type BoardRepository interface {
	BoardKey(name string) (key *BoardKey)
	DatKey(name string, parent *BoardKey) (key *DatKey)
	GetBoard(key *BoardKey, entity *BoardEntity) (err error)
	PutBoard(key *BoardKey, entity *BoardEntity) (err error)
	GetDat(key *DatKey, entity *DatEntity) (err error)
	PutDat(key *DatKey, entity *DatEntity) (err error)
	RunInTransaction(func(tx *datastore.Transaction) error) (err error)
	TxGetBoard(tx *datastore.Transaction, key *BoardKey, entity *BoardEntity) (err error)
	TxPutBoard(tx *datastore.Transaction, key *BoardKey, entity *BoardEntity) (err error)
	TxGetDat(tx *datastore.Transaction, key *DatKey, entity *DatEntity) (err error)
	TxPutDat(tx *datastore.Transaction, key *DatKey, entity *DatEntity) (err error)
}

type BoardStore struct {
	Context context.Context
	Client  *datastore.Client
}

func (repo *BoardStore) BoardKey(name string) (key *BoardKey) {
	k := datastore.NameKey(KIND_BOARD, name, nil)
	key = &BoardKey{Key: k}
	return
}

func (repo *BoardStore) DatKey(name string, parent *BoardKey) (key *DatKey) {
	k := datastore.NameKey(KIND_DAT, name, parent.Key)
	key = &DatKey{Key: k}
	return
}

func (repo *BoardStore) GetBoard(key *BoardKey, entity *BoardEntity) (err error) {
	err = repo.Client.Get(repo.Context, key.Key, entity)
	return
}

func (repo *BoardStore) PutBoard(key *BoardKey, entity *BoardEntity) (err error) {
	_, err = repo.Client.Put(repo.Context, key.Key, entity)
	return
}

func (repo *BoardStore) GetDat(key *DatKey, entity *DatEntity) (err error) {
	err = repo.Client.Get(repo.Context, key.Key, entity)
	return
}

func (repo *BoardStore) PutDat(key *DatKey, entity *DatEntity) (err error) {
	_, err = repo.Client.Put(repo.Context, key.Key, entity)
	return
}

func (repo *BoardStore) RunInTransaction(f func(tx *datastore.Transaction) error) (err error) {
	_, err = repo.Client.RunInTransaction(repo.Context, f)
	return
}

func (repo *BoardStore) TxGetBoard(tx *datastore.Transaction, key *BoardKey, entity *BoardEntity) (err error) {
	err = tx.Get(key.Key, entity)
	return
}

func (repo *BoardStore) TxPutBoard(tx *datastore.Transaction, key *BoardKey, entity *BoardEntity) (err error) {
	_, err = tx.Put(key.Key, entity)
	return
}

func (repo *BoardStore) TxGetDat(tx *datastore.Transaction, key *DatKey, entity *DatEntity) (err error) {
	err = tx.Get(key.Key, entity)
	return
}

func (repo *BoardStore) TxPutDat(tx *datastore.Transaction, key *DatKey, entity *DatEntity) (err error) {
	_, err = tx.Put(key.Key, entity)
	return
}
