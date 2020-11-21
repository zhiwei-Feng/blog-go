package controller

import (
	"github.com/kataras/iris/v12"
	"net/http"
	"path/filepath"
)

func UploadImage(ctx iris.Context) {
	_, info, err := ctx.FormFile("image")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}
	dest := filepath.Join("./uploads", info.Filename)
	ctx.Application().Logger().Debugf("dest:", dest)
	_, err = ctx.SaveFormFile(info, dest)

	var resp RespBean
	if err != nil {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "上传失败",
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "uploads/" + info.Filename,
		}
	}

	ctx.JSON(resp)
}
