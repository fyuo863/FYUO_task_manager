package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) { //Context允许在中间件之间传递变量，此处传入储存在gin中的变量
		start := time.Now() //开始计时
		ctx.Next()          //先运行下一个中间件
		duration := time.Since(start)
		fmt.Println("耗时: ", duration)
	}
}
