package middleware

import (
	"github.com/gin-gonic/gin"
)

func Catch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		urlPath := ctx.Request.URL.Path
		// ctx.JSON(http.StatusOK, gin.H{
		// 	"Method": method,
		// })

		ctx.Set("Method", method)
		ctx.Set("urlPath", urlPath)
		//log.Logger.Error("[Recovery]", "panic", method)
		ctx.Next()
	}
}
