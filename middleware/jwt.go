package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := core.New(c)
		defer ctx.Release()

		if ctx.Jwt().IsWhiteList(ctx.Request.Method, ctx.FullPath()) {
			return
		}

		// 解析 token
		md, err := ctx.Jwt().Parse()
		if md == nil {
			ctx.RespError(errors.TokenEmptyError)
			ctx.Abort()
			return
		}

		if err != nil {
			// 是否验证通过
			if err.IsVerify() {
				ctx.RespError(errors.TokenValidateError)
				ctx.Abort()
				return
			}

			// 是否过期
			if err.IsExpired() {
				ctx.RespError(errors.TokenExpiredError)
				ctx.Abort()
				return
			}

			// 判断多设备登录
			if !ctx.Jwt().CheckUnique(md.UserID) {
				ctx.RespError(errors.DulDeviceLoginError)
				ctx.Abort()
				return
			}
		}

		// 设置到上下文
		ctx.SetMetadata(md)
	}
}
