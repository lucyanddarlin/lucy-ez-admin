package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AllMenu 获取所有菜单
func AllMenu(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if resp, err := service.AllMenu(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

// AddMenu 添加菜单
func AddMenu(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.AddMenuRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.AddMenu(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}
