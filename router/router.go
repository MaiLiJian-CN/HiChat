package router

import (
	"HiChat/middleware"
	"HiChat/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("v1")
	user := v1.Group("user")
	{
		// user.GET("/list",service.List)
		user.GET("/list", service.List)
		user.POST("/login_pw", service.LoginByNameAndPassWord)
		user.POST("/new", service.NewUser)
		user.DELETE("/delete", service.DeleteUser)
		user.POST("/update", service.UpdateUser)
		user.GET("/SendUserMsg", middleware.JWY(), service.SendUserMsg)
	}
	relation := v1.Group("relation").Use(middleware.JWY())
	{
		relation.POST("/list", service.FriendList)
		relation.POST("/add", service.AddFriendByName)
		relation.POST("/new_group", service.NewGroup)
		relation.POST("/group_list", service.GroupList)
		relation.POST("/join_group", service.JoinGroup)
	}
	//聊天记录
	v1.POST("/user/redisMsg", service.RedisMsg).Use(middleware.JWY())
	return router
}
