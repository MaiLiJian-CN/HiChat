package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/middleware"
	"HiChat/models"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"

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

func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Request.FormValue("name")
	passWord := ctx.Request.FormValue("password")
	repassword := ctx.Request.FormValue("Identity")
	if user.Name == "" || passWord == "" || repassword == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "name or pwd nil",
			"data":    user,
		})
		return
	}

	//find user
	_, err := dao.FindUser(user.Name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": "user has been registered",
			"data":    user,
		})
		return
	}
	if passWord != repassword {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": " password and repwd error",
			"data":    user,
		})
		return
	}
	salt := fmt.Sprintf("%d", rand.Int31())

	//secret
	user.PassWord = common.SaltPassWord(passWord, salt)
	user.Salt = salt
	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t
	dao.CreateUser(user)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "register user seccess",
		"data":    user,
	})
}

func UpdateUser(ctx *gin.Context) {
	user := models.UserBasic{}
	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("String to Int Err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "user update error",
		})
		return
	}
	user.ID = uint(id)
	Name := ctx.Request.FormValue("name")
	PassWord := ctx.Request.FormValue("password")
	Email := ctx.Request.FormValue("email")
	Phone := ctx.Request.FormValue("phone")
	avatar := ctx.Request.FormValue("icon")
	gender := ctx.Request.FormValue("gender")

	if Email != "" {
		user.Email = Email
	}
	if Phone != "" {
		user.Phone = Phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}
	if Name != "" {
		user.Name = Name
	}
	if PassWord != "" {
		salt := fmt.Sprintf("%d", rand.Int31())
		user.Salt = salt
		user.PassWord = common.SaltPassWord(PassWord, salt)
	}


 	_, err = govalidator.ValidateStruct(user)
	if err!=nil{
		zap.S().Info("args error",err)
		ctx.JSON(http.StatusBadRequest,gin.H{
			"code":-1,
			"message":"args error",
		})

		return
	}

	rep,err:=dao.UpdateUser(user)
	if err!=nil{
		zap.S().Info("update error",err)
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"code":-1,
			"message":"update error",
		})
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"code":0,
		"message":"update success",
		"data":rep.Name,
	})
}

func DeleteUser(ctx *gin.Context){
	user:=models.UserBasic{}
	id,err:=strconv.Atoi(ctx.Request.FormValue("id"))
	if err!=nil{
		zap.S().Info("type change error",err)
		ctx.JSON(http.StatusInternalServerError,gin.H{
			"code":-1,
			"message":"delete user error",
		})
		return
	}
	user.ID=uint(id)
	err=dao.DeleteUser(user)
	if err!=nil{
		        zap.S().Info("注销用户失败", err)
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "code":    -1, //  0成功   -1失败
            "message": "注销账号失败",
        })
        return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"code":0,
		"message":"delete user success",
	})
}

