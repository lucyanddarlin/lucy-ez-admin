package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/service"
)

func Config(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	ctx.RespData(service.Config(ctx))
}
