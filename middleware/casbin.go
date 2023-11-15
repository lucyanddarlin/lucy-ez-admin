package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/model"
)

func Enforcer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := core.New(c)
		defer ctx.Release()

		path := ctx.FullPath()
		method := ctx.Request.Method

		// 获取元数据, 不存在则跳过鉴权
		md := ctx.Metadata()
		if md == nil {
			return
		}

		// 超级管理员放行
		if md.RoleKey == constants.JwtSuperAdmin {
			return
		}

		//  基础 API 放行
		menu := model.Menu{}
		if menu.IsBaseApiPath(ctx, method, path) {
			return
		}

		// 权限判断
		if is, _ := ctx.Enforcer().Instance().Enforce(md.RoleKey, path, method); !is {
			ctx.RespError(errors.NotResourcePowerError)
			ctx.Abort()
			return
		}
	}
}
