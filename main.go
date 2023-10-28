package main

import "github.com/gin-gonic/gin"

func main() {
	// 创建 gin 引擎
	engine := gin.New()

	engine.Static("/static", "./static")
}
