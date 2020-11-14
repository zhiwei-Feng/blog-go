package domain

import "time"

type Category struct {
	ID       uint `gorm:"primarykey"`
	CateName string
	Date     time.Time
}

func (Category) TableName() string {
	return "category"
}
