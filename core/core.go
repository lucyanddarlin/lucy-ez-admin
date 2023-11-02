package core

import (
	"fmt"

	"github.com/lucyanddarlin/lucy-ez-admin/config"
	"github.com/lucyanddarlin/lucy-ez-admin/core/captcha"
	"github.com/lucyanddarlin/lucy-ez-admin/core/cert"
	"github.com/lucyanddarlin/lucy-ez-admin/core/email"
	logger "github.com/lucyanddarlin/lucy-ez-admin/core/log"
	"github.com/lucyanddarlin/lucy-ez-admin/core/orm"
	"github.com/lucyanddarlin/lucy-ez-admin/core/redis"
)

func Init() {
	// 初始化配置实例
	conf := config.New()

	// 初始化全局配置
	initInstance(conf)

	// 监听配置变更要重新初始化实例
	conf.Watch(func(c *config.Config) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("配置变更失败: %v", err)
			}
		}()
		initInstance(c)
	})
}

func initInstance(conf *config.Config) {
	// 日志
	loggerIns := logger.New(conf.Log, conf.Service.Name)
	// 数据库
	ormIns := orm.New(conf.Orm, loggerIns)
	// redis
	redisIns := redis.New(conf.Redis)
	// 邮箱
	emailIns := email.New(conf.Email)
	// 验证码
	captchaIns := captcha.New(conf.Captcha, redisIns, emailIns)
	// 证书
	certIns := cert.New(conf.Cert)

	// 实例化到全局对象
	initGlobal(conf,
		withLogger(loggerIns),
		withOrm(ormIns),
		WithRedis(redisIns),
		WithEmail(emailIns),
		WithCaptcha(captchaIns),
		WithCert(certIns),
	)
}
