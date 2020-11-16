package domain

type Role struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name"`
}
