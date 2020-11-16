package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"net/http"
)

func CurrentUserName(ctx iris.Context) {
	// 通过token解析当前用户的nickname
	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)
	user := domain.User{}
	config.DB.Where("username = ?", claims.Username).First(&user)
	ctx.JSON(user.Nickname)
}

func CurrentUserId(ctx iris.Context) {
	// 通过token解析当前用户的id
	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)
	user := domain.User{}
	config.DB.Where("username = ?", claims.Username).First(&user)
	ctx.JSON(user.ID)
}

func CurrentUserEmail(ctx iris.Context) {
	// 通过token解析当前用户的email
	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)
	user := domain.User{}
	config.DB.Where("username = ?", claims.Username).First(&user)
	ctx.JSON(user.Email)
}

func IsAdmin(ctx iris.Context) {
	// 查看当前的用户是否是超级管理员
	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)
	user := domain.User{}
	config.DB.Preload("Roles").Where("username = ?", claims.Username).First(&user)
	for _, r := range user.Roles {
		if r.Name == "超级管理员" {
			ctx.JSON(true)
			return
		}
	}
	ctx.JSON(false)
}

func UpdateUserEmail(ctx iris.Context) {
	newEmail := ctx.PostValue("email")
	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)
	res := config.DB.Model(&domain.User{}).Where("id = ?", claims.UserId).Update("email", newEmail)
	var resp RespBean
	if res.Error != nil {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "开启失败",
			Data:   res.Error.Error(),
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "开启成功！",
		}
	}
	ctx.JSON(resp)
}

func GetUserByNickname(ctx iris.Context) {
	nickname := ctx.URLParam("nickname")
	var users []domain.User
	if nickname == "" {
		config.DB.Preload("Roles").Find(&users)
	} else {
		config.DB.Preload("Roles").Where("nickname = ?", nickname).Find(&users)
	}
	ctx.JSON(users)
}
