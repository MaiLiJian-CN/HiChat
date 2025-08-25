package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/models"
	"HiChat/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1, //0 success,-1 error
			"message": "Get User List Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func LoginByNameAndPassWord(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	data, err := dao.FIndUserByName(name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "login error",
		})
		return
	}
	if data.Name == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "User Not Found",
		})
		return
	}

	//由于数据库密码保存是使用md5密文的， 所以验证密码时，是将密码再次加密，然后进行对比，后期会讲解md:common.CheckPassWord
	ok := common.CheckPassWord(password, data.Salt, data.PassWord)
	if !ok {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "密码错误",
		})
		return
	}

	Rsp, err := dao.FindNuserByNameAndPwd(name, data.PassWord)
	if err != nil {
		zap.S().Info("登录失败", err)
	}

	//这里使用jwt做权限认证，后面将会介绍
	token, err := middleware.GenerateToken(Rsp.ID, "yk")
	if err != nil {
		zap.S().Info("生成token失败", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"tokens":  token,
		"userId":  Rsp.ID,
	})

}

func NewUser(ctx *gin.Context){
	user:=models.UserBasic{}
	
}
