package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"strings"
	"time"
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

type ResultForGetArticleByState struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	EditTime time.Time `json:"editTime" gorm:"column:editTime"`
	PageView int       `json:"pageView" gorm:"column:pageView"`
	State    int       `json:"state"`
	Nickname string    `json:"nickname"`
	CateName string    `json:"cateName" gorm:"column:cateName"`
}

func GetArticleByState(ctx iris.Context) {
	state := ctx.URLParamIntDefault("state", -1)
	page := ctx.URLParamIntDefault("page", -1)
	count := ctx.URLParamIntDefault("count", 6)
	keywords := ctx.URLParam("keywords")
	start := (page - 1) * count

	var claims config.UserClaims
	_ = jwt.ReadClaims(ctx, &claims)

	builder := strings.Builder{}
	builder.WriteString("%")
	builder.WriteString(keywords)
	builder.WriteString("%")

	// 1. get article count by state
	var totalCount int64
	sql := config.DB.Table("article")
	if state != -1 {
		sql = sql.Where("state = ?", state)
	}
	sql = sql.Where("uid = ?", claims.UserId)
	if keywords != "" {
		sql.Where("title LIKE ?", builder.String())
	}
	sql.Count(&totalCount)
	// 2. get article by state
	var articles []ResultForGetArticleByState
	getArticle := config.DB.Table("article a")
	getArticle = getArticle.Select("a.id, a.title, a.editTime, a.pageView, a.state, u.nickname, c.cateName")
	getArticle = getArticle.Joins("JOIN user u ON u.id = a.uid")
	getArticle = getArticle.Joins("JOIN category c ON c.id = a.cid")
	if state != -2 {
		getArticle = getArticle.Where("a.uid = ?", claims.UserId)
	}
	if state != -1 && state != -2 {
		getArticle = getArticle.Where("a.state = ?", state)
	}
	if state == -2 {
		getArticle = getArticle.Where("a.state = ?", 1)
	}
	if keywords != "" {
		getArticle = getArticle.Where("title LIKE ?", builder.String())
	}
	getArticle.Order("a.editTime desc").Limit(count).Offset(start).Find(&articles)

	// result build
	respMap := map[string]interface{}{}
	respMap["totalCount"] = totalCount
	respMap["articles"] = articles

	ctx.JSON(respMap)
}

func GetArticleByStateForAdmin(ctx iris.Context) {
	page := ctx.URLParamIntDefault("page", -1)
	count := ctx.URLParamIntDefault("count", 6)
	keywords := ctx.URLParam("keywords")
	start := (page - 1) * count

	builder := strings.Builder{}
	builder.WriteString("%")
	builder.WriteString(keywords)
	builder.WriteString("%")

	// 1. get article count by state
	var totalCount int64
	sql := config.DB.Table("article")
	sql = sql.Where("state = ?", 1)
	if keywords != "" {
		sql.Where("title LIKE ?", builder.String())
	}
	sql.Count(&totalCount)
	// 2. get article by state
	var articles []ResultForGetArticleByState
	getArticle := config.DB.Table("article a")
	getArticle = getArticle.Select("a.id, a.title, a.editTime, a.pageView, a.state, u.nickname, c.cateName")
	getArticle = getArticle.Joins("JOIN user u ON u.id = a.uid")
	getArticle = getArticle.Joins("JOIN category c ON c.id = a.cid")
	getArticle = getArticle.Where("a.state = ?", 1)
	if keywords != "" {
		getArticle = getArticle.Where("title LIKE ?", builder.String())
	}
	getArticle.Order("a.editTime desc").Limit(count).Offset(start).Find(&articles)

	// result build
	respMap := map[string]interface{}{}
	respMap["totalCount"] = totalCount
	respMap["articles"] = articles

	ctx.JSON(respMap)
}
