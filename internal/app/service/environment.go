package service

import (
	"time"
)

type BoardEnvironment interface {
	StartedAt() time.Time
	SaltComputeId() string
}

type SysEnv struct {
	StartedTime   time.Time
	ComputeIdSalt string
}

func (env *SysEnv) StartedAt() time.Time {
	return env.StartedTime
}

func (env *SysEnv) SaltComputeId() string {
	return env.ComputeIdSalt
}
