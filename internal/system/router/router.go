package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/handler"
)

func Init(engine *gin.RouterGroup) {
	api := engine.Group("/system")
	{
		// 获取系统信息
		api.GET("/config", handler.Config)

		// 用户其他操作
		api.POST("/user/login", handler.UserLogin)
	}
}
