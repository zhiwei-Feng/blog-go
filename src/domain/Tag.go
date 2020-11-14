package domain

type Tags struct {
	ID      uint `gorm:"primarykey"`
	TagName string
}

func (Tags) TableName() string {
	return "tags"
}
