package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

func UserLogin(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UserLoginRequest{
		CaptchaName: "login",
	}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现

}
