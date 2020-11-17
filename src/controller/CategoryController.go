package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAllCategories(ctx iris.Context) {
	var cates []domain.Category
	config.DB.Find(&cates)
	ctx.JSON(cates)
}

func DeleteCateById(ctx iris.Context) {
	idsArray := strings.Split(ctx.Params().Get("ids"), ",")
	var ids []int
	for _, idStr := range idsArray {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	res := config.DB.Delete(&domain.Category{}, ids)

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
			Data:   res.RowsAffected,
		}
	}
	ctx.JSON(resp)
}

func AddNewCategory(ctx iris.Context) {
	var cate domain.Category
	if err := ctx.ReadBody(&cate); err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}
	cate.Date = time.Now()

	var resp RespBean
	if cate.CateName == "" {
		resp = RespBean{
			Status: http.StatusBadRequest,
			Msg:    "请输入栏目名称",
			Data:   nil,
		}
	}
	res := config.DB.Create(&cate)
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

func UpdateCategory(ctx iris.Context) {
	var cate domain.Category
	if err := ctx.ReadBody(&cate); err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}
	res := config.DB.Model(&cate).Update("cateName", cate.CateName)
	var resp RespBean
	if res.Error != nil || res.RowsAffected != 1 {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "修改失败",
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "修改成功",
		}
	}
	ctx.JSON(resp)
}
