package domain

import "time"

type Pv struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CountDate time.Time `json:"countDate" gorm:"column:countDate"`
	Pv        int       `json:"pv"`
	Uid       int       `json:"uid"`
}

func (Pv) TableName() string {
	return "pv"
}
