package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

func LoginLog(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.LoginLogRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if list, total, err := service.PageLoginLog(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespList(total, list)
	}
}
