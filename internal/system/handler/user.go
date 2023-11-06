package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	// 创建参数 struct, 定义 captchaName 对应模板 login
	in := types.UserLoginRequest{
		CaptchaName: "login",
	}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if resp, err := service.UserLogin(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

// UserLogout 用户登出
func UserLogout(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if err := service.UserLogout(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// RefreshToken 刷新 token
func RefreshToken(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if resp, err := service.RefreshToken(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}
