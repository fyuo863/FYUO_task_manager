package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Catch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		ctx.JSON(http.StatusOK, gin.H{
			"Method": method,
		})

		ctx.Set("Method", method)
		//log.Logger.Error("[Recovery]", "panic", method)
		ctx.Next()
	}
}
