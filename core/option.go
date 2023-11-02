package core

import (
	"github.com/lucyanddarlin/lucy-ez-admin/core/captcha"
	"github.com/lucyanddarlin/lucy-ez-admin/core/cert"
	"github.com/lucyanddarlin/lucy-ez-admin/core/email"
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

func WithEmail(email email.Email) option {
	return func(g *global) {
		g.email = email
	}
}

func WithCaptcha(captcha captcha.Captcha) option {
	return func(g *global) {
		g.captcha = captcha
	}
}

func WithCert(cert cert.Cert) option {
	return func(g *global) {
		g.cert = cert
	}

}
