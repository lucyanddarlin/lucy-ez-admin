package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/constants"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/tools"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AllRole 获取全部角色
func AllRole(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if list, err := service.AllRole(ctx); err != nil {
		ctx.RespError(err)
	} else {

		ctx.RespData(list)
	}
}

// AddRole 添加角色
func AddRole(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.AddRoleRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if !tools.InList([]string{
		constants.ALLTEAM,
		constants.DOWNTEAM,
		constants.CURTEAM,
		constants.CUSTOM}, in.DataScope) {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.AddRole(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UpdateRoleRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if in.DataScope != "" && !tools.InList([]string{
		constants.ALLTEAM,
		constants.DOWNTEAM,
		constants.CURTEAM,
		constants.CUSTOM},
		in.DataScope) {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.UpdateRole(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.DeleteRoleRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.DeleteRole(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}
