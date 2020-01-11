package testutil

import (
	"cloud.google.com/go/datastore"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/dat"
	"time"
)

// A injection for google datastore
type BoardStub struct {
	BoardMap map[string]*board.Entity
	DatMap   map[string]map[string]*dat.Entity
}

func (repo *BoardStub) BoardKey(name string) (key *board.Key) {
	k := datastore.NameKey(board.KIND, name, nil)
	key = &board.Key{DSKey: k}
	return
}

func (repo *BoardStub) DatKey(name string, parent *board.Key) (key *dat.Key) {
	k := datastore.NameKey(dat.KIND, name, parent.DSKey)
	key = &dat.Key{DSKey: k}
	return
}

func (repo *BoardStub) GetBoard(key *board.Key, entity *board.Entity) (err error) {
	if e, ok := repo.BoardMap[key.DSKey.Name]; ok {
		entity.Subjects = e.Subjects
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutBoard(key *board.Key, entity *board.Entity) (err error) {
	repo.BoardMap[key.DSKey.Name] = entity
	return
}

func (repo *BoardStub) GetDat(key *dat.Key, entity *dat.Entity) (err error) {
	if board, ok := repo.DatMap[key.DSKey.Parent.Name]; !ok {
		return datastore.ErrNoSuchEntity
	} else if e, ok := board[key.DSKey.Name]; !ok {
		return datastore.ErrNoSuchEntity
	} else {
		entity.Bytes = e.Bytes
		entity.LastModified = e.LastModified
		return
	}
}

func (repo *BoardStub) PutDat(key *dat.Key, entity *dat.Entity) (err error) {
	repo.DatMap[key.DSKey.Parent.Name][key.DSKey.Name] = entity
	return
}

func (repo *BoardStub) GetAllBoard(entities *[]*board.Entity) (keys []*board.Key, err error) {
	panic("not implement")
	return
}

func (repo *BoardStub) RunInTransaction(f func(tx *datastore.Transaction) error) (err error) {
	return f(nil)
}

func (repo *BoardStub) TxGetBoard(tx *datastore.Transaction, key *board.Key, entity *board.Entity) (err error) {
	err = repo.GetBoard(key, entity)
	return
}

func (repo *BoardStub) TxPutBoard(tx *datastore.Transaction, key *board.Key, entity *board.Entity) (err error) {
	err = repo.PutBoard(key, entity)
	return
}

func (repo *BoardStub) TxGetDat(tx *datastore.Transaction, key *dat.Key, entity *dat.Entity) (err error) {
	err = repo.GetDat(key, entity)
	return
}

func (repo *BoardStub) TxPutDat(tx *datastore.Transaction, key *dat.Key, entity *dat.Entity) (err error) {
	err = repo.PutDat(key, entity)
	return
}
func (repo *BoardStub) TxGetAllBoard(tx *datastore.Transaction, entities []*board.Entity) (keys []*board.Key, err error) {
	panic("not implement")
	return
}

func (repo *BoardStub) TxPutMultiBoard(tx *datastore.Transaction, keys []*board.Key, entities []*board.Entity) (err error) {
	panic("not implement")
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
		BoardMap: make(map[string]*board.Entity),
		DatMap:   make(map[string]map[string]*dat.Entity),
	}
}

func NewBoardStub(boardName string, threads []ThreadStub) *BoardStub {
	stub := &BoardStub{
		BoardMap: map[string]*board.Entity{
			boardName: &board.Entity{
				Subjects: []board.Subject{},
			},
		},
		DatMap: map[string]map[string]*dat.Entity{
			boardName: make(map[string]*dat.Entity),
		},
	}
	boardEntity := stub.BoardMap[boardName]
	for _, v := range threads {
		boardEntity.Subjects = append(boardEntity.Subjects, board.Subject{
			ThreadKey:    v.ThreadKey,
			ThreadTitle:  v.ThreadTitle,
			MessageCount: v.MessageCount,
			LastModified: v.LastModified,
		})
		stub.DatMap[boardName][v.ThreadKey] = &dat.Entity{
			Bytes:        []byte(v.Dat),
			LastModified: v.LastModified,
		}
	}
	return stub
}
