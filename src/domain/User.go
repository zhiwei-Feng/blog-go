package domain

import (
	"time"
)

type User struct {
	ID       uint      `json:"id" gorm:"primarykey"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Nickname string    `json:"nickname"`
	Enabled  bool      `json:"enabled"`
	Roles    []Role    `json:"roles" gorm:"many2many:roles_user;foreignKey:id;joinForeignKey:uid;References:id;JoinReferences:rid"`
	Email    string    `json:"email"`
	Userface string    `json:"userface"`
	RegTime  time.Time `json:"regTime" gorm:"column:regTime"`
}

func (User) TableName() string {
	return "user"
}
