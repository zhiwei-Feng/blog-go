package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func GetAllRoles(ctx iris.Context) {
	var roles []domain.Role
	config.DB.Find(&roles)
	ctx.JSON(roles)
}

func UpdateUserRoles(ctx iris.Context) {
	userId, err := ctx.PostValueInt("id")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
	}
	// rstrs example:["1,2,3"]
	rstrs, err := ctx.PostValues("rids")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
	}
	rstrs = strings.Split(rstrs[0], ",")

	var roleIds []int
	for _, rs := range rstrs {
		rid, err := strconv.Atoi(rs)
		if err != nil {
			ctx.StopWithError(http.StatusBadRequest, err)
		}
		roleIds = append(roleIds, rid)
	}

	//+---------开始更新
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'
		// 返回任何错误都会回滚事务）

		// 1. 删除用户的所有的角色
		// 必须执行这一句，不然会出现用户有重复角色
		config.DB.Exec("DELETE FROM roles_user WHERE uid = ?", userId)
		// 2. 为用户增加提交的响应角色
		var user domain.User
		config.DB.First(&user, userId)
		for _, rid := range roleIds {
			var r domain.Role
			res := config.DB.First(&r, rid)
			if res.Error != nil {
				//找不到对应的角色
				continue
			}
			user.Roles = append(user.Roles, r)
		}
		// 这个会将user下的角色插入到中间表roles_user中，而不是进行更新
		config.DB.Save(&user)
		// end

		// 返回 nil 提交事务
		return nil
	})
	//+---------结束更新

	var resp RespBean
	if err != nil {
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
