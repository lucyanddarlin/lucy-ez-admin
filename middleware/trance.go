package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
)

// Trance 设置链路 id
func Trance() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := core.New(c)
		defer ctx.Release()

		id := ctx.GetHeader(ctx.Config().Log.Header)
		if id == "" {
			id = uuid.New().String()
		}
		ctx.SetTranceID(id)
	}
}
