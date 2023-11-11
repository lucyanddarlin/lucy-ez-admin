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

		// 发送验证码
		api.POST("/captcha", handler.Captcha)

		// 发送邮箱验证码
		api.POST("/email/captcha", handler.EmailCaptcha)

		// 角色相关
		api.GET("/roles", handler.AllRole)
		api.POST("/role", handler.AddRole)
		api.PUT("/role", handler.UpdateRole)
		api.DELETE("/role", handler.DeleteRole)

		// 部门相关
		api.GET("/teams", handler.AllTeam)
		api.POST("/team", handler.AddTeam)
		api.PUT("/team", handler.UpdateTeam)
		api.DELETE("/team", handler.DeleteTeam)

		// 用户管理相关
		api.GET("/user", handler.CurUser)
		api.PUT("/user", handler.UpdateUser)
		api.POST("/user", handler.AddUser)

		// 用户其他操作
		api.POST("/user/login", handler.UserLogin)
		api.POST("/user/logout", handler.UserLogout)
		api.POST("/token/refresh", handler.RefreshToken)
		api.GET("/login/log", handler.LoginLog)
	}
}
