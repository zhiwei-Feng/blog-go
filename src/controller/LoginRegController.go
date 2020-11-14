package controller

import (
	"github.com/kataras/iris/v12"
)

func Login(ctx iris.Context) {
	username := ctx.PostValue("username")
	password := ctx.PostValue("password")

	ctx.Application().Logger().Infof("username:[%v], password:[%v]", username, password)
	resp := RespBean{
		Status: 200,
		Msg:    "Login Success.",
	}
	ctx.JSON(resp)
}
