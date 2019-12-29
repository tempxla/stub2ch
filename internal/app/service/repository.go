package service

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
)

type BoardRepository interface {
	BoardKey(name string) (key *board.Key)
	DatKey(name string, parent *board.Key) (key *dat.Key)
	GetBoard(key *board.Key, entity *board.Entity) (err error)
	PutBoard(key *board.Key, entity *board.Entity) (err error)
	GetDat(key *dat.Key, entity *dat.Entity) (err error)
	PutDat(key *dat.Key, entity *dat.Entity) (err error)
	RunInTransaction(func(tx *datastore.Transaction) error) (err error)
	TxGetBoard(tx *datastore.Transaction, key *board.Key, entity *board.Entity) (err error)
	TxPutBoard(tx *datastore.Transaction, key *board.Key, entity *board.Entity) (err error)
	TxGetDat(tx *datastore.Transaction, key *dat.Key, entity *dat.Entity) (err error)
	TxPutDat(tx *datastore.Transaction, key *dat.Key, entity *dat.Entity) (err error)
}

type BoardStore struct {
	Context context.Context
	Client  *datastore.Client
}

func (repo *BoardStore) BoardKey(name string) (key *board.Key) {
	k := datastore.NameKey(board.KIND, name, nil)
	key = &board.Key{DSKey: k}
	return
}

func (repo *BoardStore) DatKey(name string, parent *board.Key) (key *dat.Key) {
	k := datastore.NameKey(dat.KIND, name, parent.DSKey)
	key = &dat.Key{DSKey: k}
	return
}

func (repo *BoardStore) GetBoard(key *board.Key, entity *board.Entity) (err error) {
	err = repo.Client.Get(repo.Context, key.DSKey, entity)
	return
}

func (repo *BoardStore) PutBoard(key *board.Key, entity *board.Entity) (err error) {
	_, err = repo.Client.Put(repo.Context, key.DSKey, entity)
	return
}

func (repo *BoardStore) GetDat(key *dat.Key, entity *dat.Entity) (err error) {
	err = repo.Client.Get(repo.Context, key.DSKey, entity)
	return
}

func (repo *BoardStore) PutDat(key *dat.Key, entity *dat.Entity) (err error) {
	_, err = repo.Client.Put(repo.Context, key.DSKey, entity)
	return
}

func (repo *BoardStore) RunInTransaction(f func(tx *datastore.Transaction) error) (err error) {
	_, err = repo.Client.RunInTransaction(repo.Context, f)
	return
}

func (repo *BoardStore) TxGetBoard(tx *datastore.Transaction, key *board.Key, entity *board.Entity) (err error) {
	err = tx.Get(key.DSKey, entity)
	return
}

func (repo *BoardStore) TxPutBoard(tx *datastore.Transaction, key *board.Key, entity *board.Entity) (err error) {
	_, err = tx.Put(key.DSKey, entity)
	return
}

func (repo *BoardStore) TxGetDat(tx *datastore.Transaction, key *dat.Key, entity *dat.Entity) (err error) {
	err = tx.Get(key.DSKey, entity)
	return
}

func (repo *BoardStore) TxPutDat(tx *datastore.Transaction, key *dat.Key, entity *dat.Entity) (err error) {
	_, err = tx.Put(key.DSKey, entity)
	return
}
