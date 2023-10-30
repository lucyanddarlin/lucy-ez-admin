package core

import (
	"github.com/lucyanddarlin/lucy-ez-admin/config"
	"github.com/lucyanddarlin/lucy-ez-admin/core/captcha"
	"github.com/lucyanddarlin/lucy-ez-admin/core/email"
	logger "github.com/lucyanddarlin/lucy-ez-admin/core/log"
	"github.com/lucyanddarlin/lucy-ez-admin/core/orm"
	"github.com/lucyanddarlin/lucy-ez-admin/core/redis"
)

var (
	g = new(global)
)

type global struct {
	config  *config.Config
	logger  logger.Logger
	orm     orm.Orm
	redis   redis.Redis
	email   email.Email
	captcha captcha.Captcha
}

func initGlobal(config *config.Config, opts ...option) {
	g.config = config
	for _, opt := range opts {
		opt(g)
	}
}

func GlobalConfig() *config.Config {
	return g.config
}

func GlobalLogger() logger.Logger {
	return g.logger
}

func GlobalOrm() orm.Orm {
	return g.orm
}
