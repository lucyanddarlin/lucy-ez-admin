package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/errors"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
	"github.com/lucyanddarlin/lucy-ez-admin/types"
)

// AddNotice 新建通知
func AddNotice(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.AddNoticeRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.AddNotice(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// UpdateNotice 更新通知
func UpdateNotice(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.UpdateNoticeRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.UpdateNotice(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// DeleteNotice 删除通知
func DeleteNotice(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.DeleteNoticeRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.DeleteNotice(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}

}

// PageNotice 分页获取系统通知列表
func PageNotice(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.PageNoticeRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if list, total, err := service.PageNotice(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespList(total, list)
	}
}

// GetNotice 获取通知信息
func GetNotice(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	in := types.GetNoticeRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if data, err := service.GetNotice(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(data)
	}
}

// GetNoticeUnReadNum 获取未读通知的总数
func GetNoticeUnReadNum(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	if data, err := service.GetNoticeUnReadNum(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(data)
	}
}
