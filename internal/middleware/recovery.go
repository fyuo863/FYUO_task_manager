package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() { //再包装一层func至defer内，因为要捕获panic只能在defer中
			if re := recover(); re != nil {
				log.Println("[Recovery]: ", re)
				ctx.Abort() //终止后续中间件
			}
			ctx.Next()
		}()
	}
}
