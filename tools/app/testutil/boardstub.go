package testutil

import (
	"cloud.google.com/go/datastore"
	. "github.com/tempxla/stub2ch/internal/app/types"
	"time"
)

// A injection for google datastore
type BoardStub struct {
	BoardMap map[string]*BoardEntity
	DatMap   map[string]map[string]*DatEntity
}

func (repo *BoardStub) BoardKey(name string) (key *BoardKey) {
	k := datastore.NameKey(KIND_BOARD, name, nil)
	key = &BoardKey{Key: k}
	return
}

func (repo *BoardStub) DatKey(name string, parent *BoardKey) (key *DatKey) {
	k := datastore.NameKey(KIND_DAT, name, parent.Key)
	key = &DatKey{Key: k}
	return
}

func (repo *BoardStub) GetBoard(key *BoardKey, entity *BoardEntity) (err error) {
	if e, ok := repo.BoardMap[key.Key.Name]; ok {
		entity.Subjects = e.Subjects
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutBoard(key *BoardKey, entity *BoardEntity) (err error) {
	repo.BoardMap[key.Key.Name] = entity
	return
}

func (repo *BoardStub) GetDat(key *DatKey, entity *DatEntity) (err error) {
	if board, ok := repo.DatMap[key.Key.Parent.Name]; !ok {
		return datastore.ErrNoSuchEntity
	} else if e, ok := board[key.Key.Name]; ok {
		entity.Dat = e.Dat
		return
	} else {
		return datastore.ErrNoSuchEntity
	}
}

func (repo *BoardStub) PutDat(key *DatKey, entity *DatEntity) (err error) {
	repo.DatMap[key.Key.Parent.Name][key.Key.Name] = entity
	return
}

func (repo *BoardStub) RunInTransaction(f func(tx *datastore.Transaction) error) (err error) {
	return f(nil)
}

func (repo *BoardStub) TxGetBoard(tx *datastore.Transaction, key *BoardKey, entity *BoardEntity) (err error) {
	err = repo.GetBoard(key, entity)
	return
}

func (repo *BoardStub) TxPutBoard(tx *datastore.Transaction, key *BoardKey, entity *BoardEntity) (err error) {
	err = repo.PutBoard(key, entity)
	return
}

func (repo *BoardStub) TxGetDat(tx *datastore.Transaction, key *DatKey, entity *DatEntity) (err error) {
	err = repo.GetDat(key, entity)
	return
}

func (repo *BoardStub) TxPutDat(tx *datastore.Transaction, key *DatKey, entity *DatEntity) (err error) {
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
