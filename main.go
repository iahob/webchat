package main

import (
	"embed"

	"github.com/gin-gonic/gin"
	"webchat/connection"
)

//go:embed static
var assets embed.FS

func main() {
	engine := gin.New()
	engine.GET("/", func(context *gin.Context) {
		indexBody, _ := assets.ReadFile("static/index.html")
		context.Writer.Write(indexBody)
	})
	engine.GET("/ws", connection.Upgrade)

	// 开启协程启动connection服务管理中心
	go connection.DefaultH.Run()

	// 启动http服务
	err := engine.Run(":80")
	if err != nil {
		return
	}
}
