package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
)

// AddMenu 添加菜单
func AddMenu(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()
}
