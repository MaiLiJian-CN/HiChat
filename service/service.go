package service

import (
	"HiChat/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(ctx *gin.Context){
	list,err:=dao.GetUserList()
	if err!=nil{
		ctx.JSON(200,gin.H{
			"code":-1,//0 success,-1 error
			"message":"Get User List Error",
		})
		return
	}
	ctx.JSON(http.StatusOK,list)
}