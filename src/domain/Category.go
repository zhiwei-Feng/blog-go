package domain

import "time"

type Category struct {
	ID       uint      `json:"id" gorm:"primarykey"`
	CateName string    `json:"cateName" gorm:"column:cateName"`
	Date     time.Time `json:"date"`
}

func (Category) TableName() string {
	return "category"
}
