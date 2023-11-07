package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
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
