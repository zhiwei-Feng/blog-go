package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"crypto/md5"
	"encoding/hex"
	"github.com/kataras/iris/v12"
	"net/http"
)

func Login(ctx iris.Context) {
	username := ctx.PostValue("username")
	password := ctx.PostValue("password")

	db := config.DB
	user := domain.User{}
	db.Where("username = ?", username).First(&user)

	hashedPwd := GetMD5Hash(password)
	var resp RespBean
	if hashedPwd == user.Password && user.Enabled {
		token, err := config.Authenticate(username, user.ID)
		if err != nil {
			resp = RespBean{
				Status: http.StatusUnauthorized,
				Msg:    "Token generate Failure.",
			}
		} else {
			resp = RespBean{
				Status: http.StatusOK,
				Msg:    "Login Success.",
				Data:   token,
			}
		}
	} else {
		ctx.Application().Logger().Warn("Login Failure")
		resp = RespBean{
			Status: http.StatusUnauthorized,
			Msg:    "Login Failure.",
		}
	}
	ctx.JSON(resp)
}

func Logout(ctx iris.Context) {
	ctx.Application().Logger().Info("logout")
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
