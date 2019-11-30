package testutil

import (
	E "../entity"
	"cloud.google.com/go/datastore"
)

// A injection for google datastore
type BoardStub struct {
	BoardMap map[string]*E.BoardEntity
	DatMap   map[string]map[string]*E.DatEntity
}

func (repo *BoardStub) GetBoard(key *datastore.Key, entity *E.BoardEntity) (err error) {
	if e, ok := repo.BoardMap[key.Name]; ok {
		entity.Subjects = e.Subjects
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutBoard(key *datastore.Key, entity *E.BoardEntity) (err error) {
	repo.BoardMap[key.Name] = entity
	return
}

func (repo *BoardStub) GetDat(key *datastore.Key, entity *E.DatEntity) (err error) {
	if board, ok := repo.DatMap[key.Parent.Name]; !ok {
		return datastore.ErrNoSuchEntity
	} else if e, ok := board[key.Name]; ok {
		entity.Dat = e.Dat
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutDat(key *datastore.Key, entity *E.DatEntity) (err error) {
	repo.DatMap[key.Parent.Name][key.Name] = entity
	return
}
