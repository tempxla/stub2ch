package service

import (
	"fmt"
	"github.com/tempxla/stub2ch/internal/app/memcache"
)

const (
	memkey_session_id = "SESSION_ID"
)

type BoardMemcache interface {
	Set(item *memcache.Item) error
	Get(key string) (*memcache.Item, error)
}

type AdminFunction struct {
	mem BoardMemcache
}

func (admin *AdminFunction) VerifySessionId(sessionId string) error {
	cache, err := admin.mem.Get(memkey_session_id)
	if err != nil {
		return err
	}
	if string(cache.Value) != sessionId {
		return fmt.Errorf("invalid session id")
	}
	return nil
}

func (admin *AdminFunction) Login(sessionId string) error {
	return nil
}

func (admin *AdminFunction) Logout(sessionId string) error {
	return nil
}
