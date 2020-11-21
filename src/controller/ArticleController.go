package controller

import (
	"blog-go/src/config"
	"blog-go/src/domain"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
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
		// 获取最近七天的日期
		categories = append(categories, pv.CountDate.Format("2006-01-02"))
		// 获取最近七天的数据
		dataStatistics = append(dataStatistics, pv.Pv)
	}

	respMap := map[string]interface{}{}
	respMap["categories"] = categories
	respMap["ds"] = dataStatistics

	ctx.JSON(respMap)
}

type ArticleObject struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	MdContent   string    `json:"mdContent" gorm:"column:mdContent"`
	HtmlContent string    `json:"htmlContent" gorm:"column:htmlContent"`
	Summary     string    `json:"summary"`
	Cid         uint      `json:"cid"`
	Uid         uint      `json:"uid"`
	PublishDate time.Time `json:"publishDate" gorm:"column:publishDate"`
	EditTime    time.Time `json:"editTime" gorm:"column:editTime"`
	PageView    int       `json:"pageView" gorm:"column:pageView"`
	State       int       `json:"state"`
	Nickname    string    `json:"nickname"`
	CateName    string    `json:"cateName" gorm:"column:cateName"`
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
	var articles []ArticleObject
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
	var articles []ArticleObject
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

func GetArticleById(ctx iris.Context) {
	aid, err := ctx.Params().GetInt("aid")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}

	// 1.查找article
	var article domain.Article
	config.DB.Model(&domain.Article{}).
		// todo: 数据脱敏
		Preload("Author").Preload("Cate").Preload("Tags").
		First(&article, aid)
	// 2.pv自增
	config.DB.Table("article").Where("id = ?", aid).
		UpdateColumn("pageView", gorm.Expr("pageView + ?", 1))

	ctx.JSON(article)
}

func stripHtml(content string) string {
	content = strings.ReplaceAll(content, "<p .*?>", "")
	content = strings.ReplaceAll(content, "<br\\s*/?>", "")
	content = strings.ReplaceAll(content, "\\<.*?>", "")
	return content
}

func AddNewArticle(ctx iris.Context) {
	newArticle := domain.Article{}
	if err := ctx.ReadBody(&newArticle); err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}

	if len(newArticle.DynamicTags) != 0 {
		newArticle.DynamicTags = strings.Split(newArticle.DynamicTags[0], ",")
	}

	// process
	if newArticle.Summary == "" {
		sh := stripHtml(newArticle.HtmlContent)
		if len(sh) > 50 {
			newArticle.Summary = sh[:50]
		} else {
			newArticle.Summary = sh
		}
	}
	newArticle.PublishDate = time.Now()
	newArticle.EditTime = time.Now()
	var resp RespBean
	if newArticle.ID == 0 {
		//add operate
		//设置当前用户
		var claims config.UserClaims
		_ = jwt.ReadClaims(ctx, &claims)
		newArticle.Uid = claims.UserId
		err := config.DB.Transaction(func(tx *gorm.DB) error {
			// create article
			res := tx.Create(&newArticle)
			if res.Error != nil {
				return res.Error
			}
			var tags []domain.Tags
			var existsTags []domain.Tags
			for _, t := range newArticle.DynamicTags {
				var tag domain.Tags
				if res := tx.Where("tagName = ?", t).First(&tag); res.Error != nil {
					tags = append(tags, domain.Tags{TagName: t})
				} else {
					existsTags = append(existsTags, tag)
				}
			}
			// create tags
			if len(tags) != 0 {
				res = tx.Create(&tags)
				if res.Error != nil {
					return res.Error
				}
			}

			// assemble all tag
			for _, t := range existsTags {
				tags = append(tags, t)
			}

			if len(newArticle.DynamicTags) != len(tags) {
				return errors.New("tags set fail.")
			}

			// update article
			newArticle.Tags = tags
			res = tx.Save(&newArticle)
			if res.Error != nil {
				return res.Error
			}

			return nil
		})
		if err != nil {
			ctx.StopWithError(http.StatusInternalServerError, err)
			return
		}
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    strconv.Itoa(newArticle.ID),
		}
		ctx.JSON(resp)
	} else {
		// update operate
		var oldArticle domain.Article
		config.DB.First(&oldArticle, newArticle.ID)
		// move old value to new article
		newArticle.Uid = oldArticle.Uid
		newArticle.PageView = oldArticle.PageView
		err := config.DB.Transaction(func(tx *gorm.DB) error {
			// delete all tags of article
			res := tx.Exec("DELETE FROM article_tags WHERE aid = ?", newArticle.ID)
			if res.Error != nil {
				return res.Error
			}
			// set tags
			var tags []domain.Tags
			var existsTags []domain.Tags
			for _, t := range newArticle.DynamicTags {
				var tag domain.Tags
				if res := tx.Where("tagName = ?", t).First(&tag); res.Error != nil {
					tags = append(tags, domain.Tags{TagName: t})
				} else {
					existsTags = append(existsTags, tag)
				}
			}
			// create tags
			if len(tags) != 0 {
				res = tx.Create(&tags)
				if res.Error != nil {
					return res.Error
				}
			}

			// assemble all tag
			for _, t := range existsTags {
				tags = append(tags, t)
			}

			if len(newArticle.DynamicTags) != len(tags) {
				return errors.New("tags set fail.")
			}

			// update article
			newArticle.Tags = tags
			res = tx.Save(&newArticle)
			if res.Error != nil {
				return res.Error
			}

			return nil
		})
		if err != nil {
			ctx.StopWithError(http.StatusInternalServerError, err)
			return
		}
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    strconv.Itoa(newArticle.ID),
		}
		ctx.JSON(resp)
	}
}

func UpdateArticleState(ctx iris.Context) {
	aids, err1 := ctx.PostValues("aids")
	state, err2 := ctx.PostValueInt("state")
	if err1 != nil || len(aids) == 0 {
		ctx.StopWithError(http.StatusBadRequest, err1)
		return
	}
	if err2 != nil {
		ctx.StopWithError(http.StatusBadRequest, err2)
		return
	}
	aids = strings.Split(aids[0], ",")
	var res *gorm.DB
	var resp RespBean
	var aidArray []int
	for _, aid := range aids {
		if id, err := strconv.Atoi(aid); err != nil {
			ctx.StopWithError(http.StatusBadRequest, err1)
			return
		} else {
			aidArray = append(aidArray, id)
		}
	}
	ctx.Application().Logger().Debugf("aidArray len: %v", len(aidArray))
	if state == 2 {
		// 已经在回收站内 delete Article by id
		res = config.DB.Where("id IN ?", aidArray).Delete(domain.Article{})
	} else {
		// 放入回收站, update state to 2
		res = config.DB.Model(&domain.Article{}).Where("id IN ?", aidArray).Update("state", 2)
	}
	if res.RowsAffected != int64(len(aids)) {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "删除失败",
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "删除成功",
		}
	}
	ctx.JSON(resp)
}

func RestoreArticle(ctx iris.Context) {
	id, err := ctx.PostValueInt("articleId")
	if err != nil {
		ctx.StopWithError(http.StatusBadRequest, err)
		return
	}

	var resp RespBean
	res := config.DB.Model(domain.Article{}).Where("id = ?", id).Update("state", 1)
	if res.Error != nil || res.RowsAffected == 0 {
		resp = RespBean{
			Status: http.StatusInternalServerError,
			Msg:    "还原失败",
		}
	} else {
		resp = RespBean{
			Status: http.StatusOK,
			Msg:    "还原成功",
		}
	}
	ctx.JSON(resp)
}
