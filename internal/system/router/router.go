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

		// 菜单相关
		api.GET("/menus", handler.AllMenu)
		api.POST("/menu", handler.AddMenu)
		api.PUT("/menu", handler.UpdateMenu)
		api.DELETE("/menu", handler.DeleteMenu)

		// 角色相关
		api.GET("/roles", handler.AllRole)
		api.POST("/role", handler.AddRole)
		api.PUT("/role", handler.UpdateRole)
		api.DELETE("/role", handler.DeleteRole)

		// 角色菜单相关
		api.GET("/role/menu/ids", handler.RoleMenuIds)
		api.GET("/role/menu", handler.RoleMenu)
		api.PUT("/role/menu", handler.UpdateRoleMenu)

		// 部门相关
		api.GET("/teams", handler.AllTeam)
		api.POST("/team", handler.AddTeam)
		api.PUT("/team", handler.UpdateTeam)
		api.DELETE("/team", handler.DeleteTeam)

		// 用户管理相关
		api.GET("users", handler.PageUser)
		api.GET("/user", handler.CurUser)
		api.PUT("/user", handler.UpdateUser)
		api.POST("/user", handler.AddUser)
		api.POST("user/info", handler.UpdateUserInfo)
		api.POST("/user/verify", handler.UpdateUserInfoByVerify)
		api.DELETE("user", handler.DeleteUser)
		api.GET("/user/menus", handler.UserMenus)

		// 用户其他操作
		api.POST("/user/login", handler.UserLogin)
		api.POST("/user/logout", handler.UserLogout)
		api.POST("/token/refresh", handler.RefreshToken)
		api.GET("/login/log", handler.LoginLog)

		// 系统通知相关
		api.POST("/notice", handler.AddNotice)

	}
}
