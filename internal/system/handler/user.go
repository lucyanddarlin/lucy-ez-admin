package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
)

func UserLogin(c *gin.Context) {
	ctx := core.New(c)
	defer ctx.Release()

	ctx.RespSuccess()
}
