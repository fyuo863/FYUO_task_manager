package router

import (
	"FYUO_task_manager/internal/middleware"

	"FYUO_task_manager/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine { //Gin 运行模式
	gin.SetMode(mode)
	r := gin.New()               //创建gin路由
	r.Use(middleware.Recovery()) //gin.Recovery(),捕获所有 panic
	r.Use(middleware.Logger())   //gin.Logger(),记录每个请求的耗时和状态码
	r.Use(middleware.Catch())    //捕获请求方法

	menu := r.Group("/api")
	{
		menu.GET("/test", handler.Test)
	}

	// v1 := r.Group("/api/v1"){

	// }

	return r
}
