package service

import (
	"time"
)

type SysEnv struct {
	CurrentTime   time.Time
	ComputeIdSalt string
	AdminMailSalt string
}

func (env *SysEnv) Now() time.Time {
	return env.CurrentTime
}

func (env *SysEnv) SaltComputeId() string {
	return env.ComputeIdSalt
}

func (env *SysEnv) SaltAdminMail() string {
	return env.AdminMailSalt
}
