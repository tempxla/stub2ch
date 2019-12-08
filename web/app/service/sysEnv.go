package service

import (
	"time"
)

type SysEnv struct {
	StartedTime   time.Time
	ComputeIdSalt string
	AdminMailSalt string
}

func (env *SysEnv) StartedAt() time.Time {
	return env.StartedTime
}

func (env *SysEnv) SaltComputeId() string {
	return env.ComputeIdSalt
}

func (env *SysEnv) SaltAdminMail() string {
	return env.AdminMailSalt
}
