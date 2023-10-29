package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/install"
	"github.com/lucyanddarlin/lucy-ez-admin/internal/system/router"
	"github.com/lucyanddarlin/lucy-ez-admin/middleware"
)

func main() {
	// 创建 gin 引擎
	engine := gin.New()

	// 核心组件初始化
	core.Init()

	// 进行系统初始化
	install.Init()

	// 初始化静态资源
	engine.Static("/static", "./static")

	// 注册中间件
	api := middleware.Registry(engine)

	// 路由初始化
	router.Init(api)

	// 启动服务
	srv := core.GlobalConfig().Service
	if err := engine.Run(srv.Addr); err != nil {
		log.Fatalln(err)
	}
}
