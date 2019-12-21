package service

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/tempxla/stub2ch/configs/app/admincfg"
	"github.com/tempxla/stub2ch/internal/app/util"
	"time"
)

type AdminFunction struct {
	mem BoardMemcache
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
	item := &Item{
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

	if err := admin.mem.Delete(admincfg.LOGIN_COOKIE_NAME); err != nil {
		return err
	}

	return nil
}
