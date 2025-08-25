package middleware

import (
	"errors"
	"net/http"
	"os/user"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var(
	TokenExpired = errors.New("Token is expired")
)


// secret key

var jwtSecret=[]byte("ice_moss")

type Claims struct{
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func GenerateToken(userId uint,iss string )(string,error){
	//token valid
	nowTime :=time.Now()
	expireTime:=nowTime.Add(48*30*time.Hour)

	claims:=Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer	: iss,
		},
	}

	tokenClaims:=jwt.NewWithClaims(jwt.SigningMethodES256,claims)

	token,err:=tokenClaims.SignedString(jwtSecret)
	return token,err

}

// func JWT() gin.HandlerFunc{
// 	return func(ctx *gin.Context) {
// 		token:=ctx.PostForm("token")
// 		user:=ctx.Query("userId")
// 		userId,err:=strconv.Atoi(user)
// 		if err!=nil {
// 			c.JSON
// 		}
// 	}
// }
