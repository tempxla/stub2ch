package service

import (
	"cloud.google.com/go/datastore"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"github.com/tempxla/stub2ch/internal/app/service/repository"
	"github.com/tempxla/stub2ch/internal/app/types/entity/board"
	"github.com/tempxla/stub2ch/internal/app/types/entity/memcache"
	"github.com/tempxla/stub2ch/internal/app/util"
	"log"
	"time"
)

type AdminFunction struct {
	repo repository.BoardRepository
	mem  BoardMemcache
}

func (admin *AdminFunction) VerifySession(sessionId string) error {
	cache, err := admin.mem.Get(admincfg.LOGIN_COOKIE_NAME)
	if err != nil {
		return err
	}
	if string(cache.Value) != sessionId {
		return fmt.Errorf("invalid session id.")
	}
	return nil
}

func (admin *AdminFunction) Login(passphrase, signature string) (string, error) {

	sha256phrase := sha256.Sum256([]byte(passphrase))

	if fmt.Sprintf("%x", sha256phrase) != admincfg.LOGIN_PASSPHRASE_DIGEST {
		return "", fmt.Errorf("invalid passphrase.")
	}

	if err := util.VerifyPKCS1v15(admincfg.RSA_PUBLIC, sha256phrase, signature); err != nil {
		return "", err
	}

	sessionId := uuid.New().String()
	item := &memcache.Item{
		Key:        admincfg.LOGIN_COOKIE_NAME,
		Value:      []byte(sessionId),
		Expiration: time.Duration(30) * time.Minute, // 30分たっても削除されないよ！
	}
	err := admin.mem.Set(item)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (admin *AdminFunction) Logout() error {
	return admin.mem.Delete(admincfg.LOGIN_COOKIE_NAME)
}

// 空の板を作成する。
// すでに存在する場合エラーを返す。
func (admin *AdminFunction) CreateBoard(boardName string) error {
	log.Printf("CreateBoard: %v", boardName)

	key := admin.repo.BoardKey(boardName)
	newEntity := &board.Entity{
		Subjects: []board.Subject{},
	}
	entity := &board.Entity{}

	return admin.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		err := admin.repo.GetBoard(key, entity)
		if err != nil && err != datastore.ErrNoSuchEntity {
			// datastoreのエラー
			return err
		}
		if err == nil {
			// already exists.
			return fmt.Errorf("entity duplicated: %v", boardName)
		}
		// err == datastore.ErrNoSuchEntity
		return admin.repo.PutBoard(key, newEntity)
	})
}

func (admin *AdminFunction) GetWriteCount() (_ int, err error) {

	var entities []*board.Entity
	_, err = admin.repo.GetAllBoard(&entities)
	if err != nil {
		return
	}

	count := 0
	for _, entity := range entities {
		count += entity.WriteCount
	}

	return count, nil
}

func (admin *AdminFunction) ResetWriteCount() error {
	return admin.repo.RunInTransaction(func(tx *datastore.Transaction) error {
		var entities []*board.Entity
		keys, err := admin.repo.TxGetAllBoard(tx, &entities)
		if err != nil {
			return err
		}

		for i, _ := range keys {
			entities[i].WriteCount = 0
		}

		return admin.repo.TxPutMultiBoard(tx, keys, entities)
	})
}
