package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AllTeam 获取所有部门
func AllTeam(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if tree, err := service.AllTeam(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(tree)
	}
}

// AddTeam 添加部门
func AddTeam(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.AddTeamRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.AddTeam(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// UpdateTeam 更新部门信息
func UpdateTeam(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UpdateTeamRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.UpdateTeam(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// DeleteTeam 删除部门
func DeleteTeam(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.DeleteTeamRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if err := service.DeleteTeam(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}
