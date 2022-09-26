package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SUCCESS int = 0 // 操作成功
	FAILED  int = 1 // 操作失败
)

// Success 请求成功的时候 使用该方法返回信息
func Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": SUCCESS,
		"msg":  msg,
		"data": data,
	})
}

// Fail 请求失败的时候, 使用该方法返回信息
func Fail(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": FAILED,
		"msg":  msg,
		"data": nil,
	})
}
