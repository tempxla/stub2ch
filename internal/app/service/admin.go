package service

import (
	"fmt"
	"github.com/tempxla/stub2ch/configs/app/config"
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
