package domain

import "time"

type Article struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
	MdContent   string    `json:"mdContent" gorm:"column:mdContent"`
	HtmlContent string    `json:"htmlContent" gorm:"column:htmlContent"`
	Summary     string    `json:"summary"`
	Cid         uint      `json:"cid"`
	Uid         uint      `json:"uid"`
	PublishDate time.Time `json:"publishDate" gorm:"column:publishDate"`
	State       int       `json:"state"`
	PageView    int       `json:"pageView" gorm:"column:pageView"`
	EditTime    time.Time `json:"editTime" gorm:"column:editTime"`
	DynamicTags *[]string `json:"dynamicTags"`
	Nickname    string    `json:"nickname"`
	CateName    string    `json:"cateName" gorm:"column:cateName"`
	Tags        *[]Tags   `json:"tags"`
	StateStr    string    `json:"stateStr"`
}

func (Article) TableName() string {
	return "article"
}
