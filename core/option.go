package core

import (
	logger "github.com/lucyanddarlin/lucy-ez-admin/core/log"
	"github.com/lucyanddarlin/lucy-ez-admin/core/orm"
	"github.com/lucyanddarlin/lucy-ez-admin/core/redis"
)

type option func(*global)

func withLogger(log logger.Logger) option {
	return func(g *global) {
		g.logger = log
	}
}

func withOrm(orm orm.Orm) option {
	return func(g *global) {
		g.orm = orm
	}
}

func WithRedis(redis redis.Redis) option {
	return func(g *global) {
		g.redis = redis
	}
}
