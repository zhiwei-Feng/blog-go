package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"net/http"
	"strings"
	"time"
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
		builder := strings.Builder{}
		builder.WriteString("%")
		builder.WriteString(nickname)
		builder.WriteString("%")
		config.DB.Preload("Roles").Where("nickname Like ?", builder.String()).Find(&users)
	}
	ctx.JSON(users)
}

func GetUserById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}
	user := domain.User{}
	config.DB.Preload("Roles").First(&user, id)
	ctx.JSON(user)
}

func DeleteUserById(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}
	// 会将roles_user中对应用户的记录一并删除
	res := config.DB.Delete(&domain.User{ID: uint(id)})
	var resp RespBean
	if res.Error != nil {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "删除失败",
			Data:   nil,
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "删除成功",
			Data:   nil,
		}
	}
	ctx.JSON(resp)
}

func SwitchUserEnableStatus(ctx iris.Context) {
	enabled, err1 := ctx.PostValueBool("enabled")
	id, err2 := ctx.PostValueInt("uid")
	if err1 != nil || err2 != nil {
		ctx.StopWithStatus(http.StatusBadRequest)
		return
	}

	user := domain.User{}
	res := config.DB.First(&user, id)
	if res.Error != nil {
		ctx.StopWithError(http.StatusInternalServerError, res.Error)
		return
	}
	res = config.DB.Model(&user).Update("enabled", enabled)

	var resp RespBean
	if res.Error != nil {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "更新失败",
			Data:   nil,
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "更新成功",
			Data:   nil,
		}
	}
	ctx.JSON(resp)

}

type UserObject struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Nickname string    `json:"nickname"`
	Enabled  bool      `json:"enabled"`
	Roles    []string  `json:"roles"`
	Email    string    `json:"email"`
	Userface string    `json:"userface"`
	RegTime  time.Time `json:"regTime"`
}

func AddUser(ctx iris.Context) {
	var resp RespBean
	var newUserObj UserObject
	if err := ctx.ReadBody(&newUserObj); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}
	newUserObj.Roles = strings.Split(newUserObj.Roles[0], ",")
	var newU domain.User
	newU.Username = newUserObj.Username
	newU.Nickname = newUserObj.Nickname
	newU.Password = GetMD5Hash(newUserObj.Password)
	newU.Email = newUserObj.Email
	newU.Userface = newUserObj.Userface
	for _, idStr := range newUserObj.Roles {
		var temp domain.Role
		err := config.DB.First(&temp, idStr)
		if err.Error != nil {
			ctx.StopWithError(http.StatusBadRequest, err.Error)
			return
		}
		newU.Roles = append(newU.Roles, temp)
	}
	newU.Enabled = false
	newU.RegTime = time.Now()
	res := config.DB.Create(&newU)
	if res.Error != nil {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "添加失败",
			Data:   nil,
		}

	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "添加成功",
			Data:   nil,
		}
	}
	ctx.JSON(resp)
}
