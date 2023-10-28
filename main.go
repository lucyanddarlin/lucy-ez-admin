package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lucyanddarlin/lucy-ez-admin/core"
)

func main() {
	// 创建 gin 引擎
	engine := gin.New()

	engine.Static("/static", "./static")

	// 核心组件初始化
	core.Init()

	// 启动服务
	srv := core.GlobalConfig().Service
	if err := engine.Run(srv.Addr); err != nil {
		log.Fatalln(err)
	}
}
