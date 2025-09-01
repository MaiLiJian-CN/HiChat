package service

import (
	"fmt"
	"strconv"

	"HiChat/common"
	"HiChat/dao"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//user对返回数据进行屏蔽
type user struct {
    Name     string
    Avatar   string
    Gender   string
    Phone    string
    Email    string
    Identity string
}

func FriendList(ctx *gin.Context){
	id, _ := strconv.Atoi(ctx.Request.FormValue("userId"))
	fmt.Println("id=",id)
    users, err := dao.FriendList(uint(id))
    if err != nil {
        zap.S().Info("获取好友列表失败", err)
        ctx.JSON(200, gin.H{
            "code":    -1, //  0成功   -1失败
            "message": "好友为空",
        })
        return
    }


    infos := make([]user, 0)

    for _, v := range *users {
        info := user{
            Name:     v.Name,
            Avatar:   v.Avatar,
            Gender:   v.Gender,
            Phone:    v.Phone,
            Email:    v.Email,
            Identity: v.Identity,
        }
        infos = append(infos, info)
    }
    common.RespOKList(ctx.Writer, infos, len(infos))

}

func AddFriendByName(ctx *gin.Context){
	user:=ctx.PostForm("userId")
	userId,err:=strconv.Atoi(user)
	if 	err!=nil{
		zap.S().Info("type change error",err)
		return
	}
	tar:=ctx.PostForm("targetName")
	target,err:=strconv.Atoi(tar)
	if err!=nil{
		code,err:=dao.AddFriendByName(uint(userId),tar)
		if err!=nil{
			HandleErr(code,ctx,err)
			return
		}
	}else{
		code,err:=dao.AddFriend(uint(userId),uint(target))
		if err!=nil{
			HandleErr(code,ctx,err)
			return
		}
	}
	ctx.JSON(200,gin.H{
		"code":0,
		"msg":"add Friend success",
	})
}
func HandleErr(code int,ctx *gin.Context,err error){
	switch code{
	case -1:
		ctx.JSON(200,gin.H{
			"code":-1,
			"msg":err.Error(),
		})
	case 0:
		ctx.JSON(200,gin.H{
			"code":-1,
			"msg":"Friend is added",
		})
	case -2:
		ctx.JSON(200,gin.H{
			"code":-1,
			"msg":"can not add self",
		})
	}
}