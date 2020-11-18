package domain

import "time"

type Article struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
	MdContent   string    `json:"mdContent" gorm:"column:mdContent"`
	HtmlContent string    `json:"htmlContent" gorm:"column:htmlContent"`
	Summary     string    `json:"summary"`
	Cid         uint      `json:"cid"`
	Uid         uint      `json:"uid"`
	PublishDate time.Time `json:"publishDate" gorm:"column:publishDate"`
	EditTime    time.Time `json:"editTime" gorm:"column:editTime"`
	State       int       `json:"state"`
	PageView    int       `json:"pageView" gorm:"column:pageView"`

	DynamicTags []string `json:"dynamicTags" gorm:"-"`
	Author      User     `json:"author" gorm:"foreignKey:id;References:uid"`
	Cate        Category `json:"category" gorm:"foreignKey:id;References:cid"`
	Tags        []Tags   `json:"tags" gorm:"many2many:article_tags;foreignKey:id;joinForeignKey:aid;References:id;JoinReferences:tid"`
	//StateStr string   `json:"stateStr"`
}

func (Article) TableName() string {
	return "article"
}
