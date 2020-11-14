package domain

import "time"

type Article struct {
	ID          uint `gorm:"primarykey"`
	Title       string
	MdContent   string
	HtmlContent string
	Summary     string
	Cid         uint
	Uid         uint
	PublishDate time.Time
	State       int
	PageView    int
	EditTime    time.Time
	DynamicTags []string
	Nickname    string
	CateName    string
	Tags        []Tags
	StateStr    string
}

func (Article) TableName() string {
	return "article"
}
