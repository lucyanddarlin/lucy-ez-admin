package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
)

func Registry(engine *gin.Engine) *gin.RouterGroup {
	conf := core.GlobalConfig()

	// 开启全局 404
	engine.NoRoute(Resp404())

	// 开启健康检查
	engine.GET("/healthy", Healthy())

	// 开启 CORS
	if conf.Middleware.Cors.Enable {
		engine.Use(Cors())
	}

	// 开启链路 id 和全局异常捕捉异常恢复
	engine.Use(Trance(), Recovery())

	api := engine.Group("api")

	// 开启请求日志
	api.Use(RequestLog())

	// 开启 jwt 鉴权
	api.Use(JwtAuth())

	return api

}
