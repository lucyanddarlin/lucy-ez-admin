package config

import "time"

type Orm struct {
	Enable        string
	Drive         string
	Name          string
	Dsn           string
	MaxLifeTime   time.Duration
	MaxOpenConn   int
	MaxIdleConn   int
	Level         int
	SlowThreshold time.Duration
	Replicas      []string
}
