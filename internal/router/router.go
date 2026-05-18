package router

import (
	"FYUO_task_manager/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine { //Gin 运行模式
	gin.SetMode(mode)
	r := gin.New()               //创建gin路由
	r.Use(middleware.Recovery()) //gin.Recovery(),捕获所有 panic
	r.Use(middleware.Logger())   //gin.Logger(),记录每个请求的耗时和状态码

	return r
}
