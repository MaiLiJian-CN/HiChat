package router

import (
	"HiChat/service"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine{
	router:=gin.Default()

	v1:=router.Group("v1")
	user:=v1.Group("user")
	{
		// user.GET("/list",service.List)	
		user.GET("/list",service.List)
		user.POST("/login_pw",service.LoginByNameAndPassWord)
		user.POST("/new",service.NewUser)
		user.DELETE("/delete",service.DeleteUser)
		user.POST("/update",service.UpdateUser)
	}
	return router
}
