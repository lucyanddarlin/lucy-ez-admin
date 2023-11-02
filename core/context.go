package core

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/config"
	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core/captcha"
	"github.com/lucyanddarlin/lucy-ez-admin/core/cert"
	"github.com/lucyanddarlin/lucy-ez-admin/core/http"
	"github.com/lucyanddarlin/lucy-ez-admin/core/jwt"
	"github.com/lucyanddarlin/lucy-ez-admin/core/orm"
	"github.com/lucyanddarlin/lucy-ez-admin/core/redis"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
	"go.uber.org/zap"
)

var ctxPool = sync.Pool{
	New: func() any {
		return &Context{}
	},
}

type Context struct {
	*gin.Context
}

func New(ctx *gin.Context) *Context {
	c := ctxPool.Get().(*Context)
	c.Context = ctx
	return c
}

// Gin 返回 gin 上下文
func (ctx *Context) Gin() *gin.Context {
	return ctx.Context
}

// Release 释放 ctx 到 pool 中
func (ctx *Context) Release() {
	ctx.Context = nil
	ctxPool.Put(ctx)
}

// Config 获取配置文件信息
func (ctx *Context) Config() *config.Config {
	return g.config
}

// Cert 获取证书实例
func (ctx *Context) Cert() cert.Cert {
	return g.cert
}

// TranceID 获取链路日志 ID
func (ctx *Context) TranceID() string {
	return ctx.Gin().GetString(ctx.Config().Log.Field)
}

// SetTranceID 设置链路日志 ID
func (ctx *Context) SetTranceID(id string) {
	ctx.Gin().Set(ctx.Config().Log.Field, id)
}

// Logger 获取文件日志器
func (ctx *Context) Logger() *zap.Logger {
	return g.logger.WithID(ctx.TranceID())
}

// Orm 获取数据库
func (ctx *Context) Orm() orm.Orm {
	return g.orm
}

// Redis 获取 Redis
func (ctx *Context) Redis() redis.Redis {
	return g.redis
}

// Http 获取请求器
func (ctx *Context) Http() http.Request {
	return http.New(ctx.Config().Http, ctx.Logger())
}

func (ctx *Context) Metadata() *types.Metadata {
	val, is := ctx.Get(constants.Metadata)
	if !is {
		return nil
	}
	meta, is := val.(*types.Metadata)
	if !is {
		return nil
	}
	return meta
}

func (ctx *Context) SetMetadata(val *types.Metadata) {
	ctx.Set(constants.Metadata, val)
}

func (ctx *Context) SourceCtx() context.Context {
	c := context.Background()
	for key, val := range ctx.Context.Keys {
		c = context.WithValue(c, key, val)
	}
	return c
}

// ImageCaptcha 生成图片验证码实例
//
//	@param name 验证码模板名称
func (ctx *Context) ImageCaptcha(name string) captcha.Image {
	return g.captcha.Image(ctx.ClientIP(), name)
}

func (ctx *Context) ClientIP() string {
	ip := ctx.Context.ClientIP()
	if ip == "::1" {
		ip = ctx.GetHeader("X-Real-IP")
	}
	return ip
}

func (ctx *Context) Jwt() jwt.JWT {
	conf := ctx.Config().JWT
	token := ctx.Gin().GetHeader(conf.Header)
	return jwt.New(conf, ctx.Redis(), token)
}
