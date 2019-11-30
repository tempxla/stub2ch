package service

import (
	E "../entity"
	"cloud.google.com/go/datastore"
	"context"
)

type BoardStore struct {
	Context context.Context
	Client  *datastore.Client
}

func (repo *BoardStore) GetBoard(key *datastore.Key, entity *E.BoardEntity) (err error) {
	if err = repo.Client.Get(repo.Context, key, entity); err != nil {
		return
	}
	return
}

func (repo *BoardStore) PutBoard(key *datastore.Key, entity *E.BoardEntity) (err error) {
	if _, err := repo.Client.Put(repo.Context, key, entity); err != nil {
		return err
	}
	return
}

func (repo *BoardStore) GetDat(key *datastore.Key, entity *E.DatEntity) (err error) {
	if err = repo.Client.Get(repo.Context, key, entity); err != nil {
		return
	}
	return
}

func (repo *BoardStore) PutDat(key *datastore.Key, entity *E.DatEntity) (err error) {
	if _, err := repo.Client.Put(repo.Context, key, entity); err != nil {
		return err
	}
	return
}
