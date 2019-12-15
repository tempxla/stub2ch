package service

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"github.com/tempxla/stub2ch/configs/app/config"
	"github.com/tempxla/stub2ch/internal/app/util"
	"time"
)

type AdminFunction struct {
	mem BoardMemcache
}

func (admin *AdminFunction) VerifySessionId(sessionId string) error {
	cache, err := admin.mem.Get(config.ADMIN_COOKIE_NAME)
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

	if string(sha256phrase[:]) != config.ADMIN_PASSPHRASE_DIGEST {
		return "", fmt.Errorf("invalid passphrase.")
	}

	if err := util.VerifyPKCS1v15(config.RSA_PUBLIC, sha256phrase, signature); err != nil {
		return "", err
	}

	sessionId := uuid.New().String()
	item := &Item{
		Key:        config.ADMIN_COOKIE_NAME,
		Value:      []byte(sessionId),
		Expiration: time.Duration(30) * time.Minute, // 30分たっても削除されないよ！
	}
	err := admin.mem.Set(item)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (admin *AdminFunction) Logout(sessionId string) error {

	if err := admin.mem.Delete(config.ADMIN_COOKIE_NAME); err != nil {
		return err
	}

	return nil
}
