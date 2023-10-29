package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
)

func Resp404() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := core.New(c)
		defer ctx.Release()

		ctx.RespError(errors.New("将要处理的接口不存在"))
	}
}
