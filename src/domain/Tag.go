package domain

type Tags struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	TagName string `json:"tagName" gorm:"column:tagName"`
}

func (Tags) TableName() string {
	return "tags"
}
