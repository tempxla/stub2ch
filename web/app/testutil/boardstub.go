package testutil

import (
	. "../entity"
	"cloud.google.com/go/datastore"
	"time"
)

// A injection for google datastore
type BoardStub struct {
	BoardMap map[string]*BoardEntity
	DatMap   map[string]map[string]*DatEntity
}

func (repo *BoardStub) GetBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	if e, ok := repo.BoardMap[key.Name]; ok {
		entity.Subjects = e.Subjects
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutBoard(key *datastore.Key, entity *BoardEntity) (err error) {
	repo.BoardMap[key.Name] = entity
	return
}

func (repo *BoardStub) GetDat(key *datastore.Key, entity *DatEntity) (err error) {
	if board, ok := repo.DatMap[key.Parent.Name]; !ok {
		return datastore.ErrNoSuchEntity
	} else if e, ok := board[key.Name]; ok {
		entity.Dat = e.Dat
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutDat(key *datastore.Key, entity *DatEntity) (err error) {
	repo.DatMap[key.Parent.Name][key.Name] = entity
	return
}

func (repo *BoardStub) RunInTransaction(f func(tx *datastore.Transaction) error) (err error) {
	return f(nil)
}

func (repo *BoardStub) TxGetBoard(tx *datastore.Transaction, key *datastore.Key, entity *BoardEntity) (err error) {
	err = repo.GetBoard(key, entity)
	return
}

func (repo *BoardStub) TxPutBoard(tx *datastore.Transaction, key *datastore.Key, entity *BoardEntity) (err error) {
	err = repo.PutBoard(key, entity)
	return
}

func (repo *BoardStub) TxGetDat(tx *datastore.Transaction, key *datastore.Key, entity *DatEntity) (err error) {
	err = repo.GetDat(key, entity)
	return
}

func (repo *BoardStub) TxPutDat(tx *datastore.Transaction, key *datastore.Key, entity *DatEntity) (err error) {
	err = repo.PutDat(key, entity)
	return
}

type ThreadStub struct {
	ThreadKey    string
	ThreadTitle  string
	MessageCount int
	LastModified time.Time
	Dat          string
}

func EmptyBoardStub() *BoardStub {
	return &BoardStub{
		BoardMap: make(map[string]*BoardEntity),
		DatMap:   make(map[string]map[string]*DatEntity),
	}
}

func NewBoardStub(boardName string, threads []ThreadStub) *BoardStub {
	stub := &BoardStub{
		BoardMap: map[string]*BoardEntity{
			boardName: &BoardEntity{
				Subjects: []Subject{},
			},
		},
		DatMap: map[string]map[string]*DatEntity{
			boardName: make(map[string]*DatEntity),
		},
	}
	board := stub.BoardMap[boardName]
	for _, v := range threads {
		board.Subjects = append(board.Subjects, Subject{
			ThreadKey:    v.ThreadKey,
			ThreadTitle:  v.ThreadTitle,
			MessageCount: v.MessageCount,
			LastModified: v.LastModified,
		})
		stub.DatMap[boardName][v.ThreadKey] = &DatEntity{Dat: []byte(v.Dat)}
	}
	return stub
}
