package service

import (
	"HiChat/models"

	"github.com/gin-gonic/gin"
)

//SendUserMsg 发送消息
func SendUserMsg(ctx *gin.Context) {
    models.Chat(ctx.Writer, ctx.Request)
}