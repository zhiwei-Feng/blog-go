package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
)

func CurrentUserName(ctx iris.Context) {
	// 通过token解析当前用户的nickname
	var claims config.UserClaims
	if err := config.J.VerifyToken(ctx, &claims); err != nil {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}
	user := domain.User{}
	config.DB.Where("username = ?", claims.Username).First(&user)
	ctx.JSON(user.Nickname)
}

func CurrentUserId(ctx iris.Context) {
	// 通过token解析当前用户的id
	var claims config.UserClaims
	if err := config.J.VerifyToken(ctx, &claims); err != nil {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}
	user := domain.User{}
	config.DB.Where("username = ?", claims.Username).First(&user)
	ctx.JSON(user.ID)
}
