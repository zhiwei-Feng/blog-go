package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func GetDataStatistics(ctx iris.Context) {
	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)

	var pvs []domain.Pv
	config.DB.Where("uid = ?", claims.UserId).Order("countDate").Limit(7).Find(&pvs)

	var (
		categories     []string
		dataStatistics []int
	)

	for _, pv := range pvs {
		categories = append(categories, pv.CountDate.Format("2006-01-02"))
		dataStatistics = append(dataStatistics, pv.Pv)
	}

	respMap := map[string]interface{}{}
	respMap["categories"] = categories
	respMap["ds"] = dataStatistics

	ctx.JSON(respMap)
}
