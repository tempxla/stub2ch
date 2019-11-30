package services

import (
	"cloud.google.com/go/datastore"
	"context"
)

type BoardStore struct {
	ctx    context.Context
	client *datastore.Client
}

func (repo *BoardStore) GetBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	if err = repo.client.Get(repo.ctx, key, entity); err != nil {
		return
	}
	return
}

func (repo *BoardStore) PutBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	if _, err := repo.client.Put(repo.ctx, key, entity); err != nil {
		return err
	}
	return
}

func (repo *BoardStore) GetDat(key *datastore.Key, entity *DatEntity) (err error) {
	if err = repo.client.Get(repo.ctx, key, entity); err != nil {
		return
	}
	return
}

func (repo *BoardStore) PutDat(key *datastore.Key, entity *DatEntity) (err error) {
	if _, err := repo.client.Put(repo.ctx, key, entity); err != nil {
		return err
	}
	return
}
