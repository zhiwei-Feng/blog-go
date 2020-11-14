package controller

import "github.com/kataras/iris/v12"

func CurrentUserName(ctx iris.Context) {
	// 通过token解析当前用户的nickname
	ctx.JSON("admin")
}

func CurrentUserId(ctx iris.Context) {
	// 通过token解析当前用户的id
	ctx.JSON(6)
}


