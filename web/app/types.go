package main

import (
	"cloud.google.com/go/datastore"
	"context"
	"time"
)

// Kind=Board
// Key=BoardName
type BoardEntity struct {
	Subjects Subjects `datastore:",noindex"`
}

type Subjects []Subject

type Subject struct {
	ThreadKey    string
	ThreadTitle  string
	MessageCount int
	LastFloat    time.Time
	LastModified time.Time
}

func (s Subjects) Len() int {
	return len(s)
}

func (s Subjects) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Subjects) Less(i, j int) bool {
	return s[i].LastFloat.Before(s[j].LastFloat)
}

// Kind=Dat
// Ancestor=Board
// Key=ThreadKey
type DatEntity struct {
	Dat []byte
}

// Dependency injection for Board
type BoardService struct {
	repo BoardRepository
}

func NewBoardService(repo BoardRepository) *BoardService {
	return &BoardService{
		repo: repo,
	}
}

type BoardRepository interface {
	GetBoard(key *datastore.Key, entity *BoardEntity) (err error)
	PutBoard(key *datastore.Key, entity *BoardEntity) (err error)
	GetDat(key *datastore.Key, entity *DatEntity) (err error)
	PutDat(key *datastore.Key, entity *DatEntity) (err error)
}

func (sv *BoardService) GetBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	return sv.repo.GetBoard(key, entity)
}
func (sv *BoardService) PutBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	return sv.repo.PutBoard(key, entity)
}
func (sv *BoardService) GetDat(key *datastore.Key, entity *DatEntity) (err error) {
	return sv.repo.GetDat(key, entity)
}
func (sv *BoardService) PutDat(key *datastore.Key, entity *DatEntity) (err error) {
	return sv.repo.PutDat(key, entity)
}

// A injection for google datastore
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
