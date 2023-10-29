package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
	"github.com/lucyanddarlin/lucy-ez-admin/install"
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

	// 启动服务
	srv := core.GlobalConfig().Service
	if err := engine.Run(srv.Addr); err != nil {
		log.Fatalln(err)
	}
}
