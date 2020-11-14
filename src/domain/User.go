package domain

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	Username  string
	Password  string
	Nickname  string
	Enabled   bool
	Roles     []Role `gorm:"many2many:roles_user;foreignKey:id;joinForeignKey:uid;References:id;JoinReferences:rid"`
	Email     string
	Userface  string
	Timestamp time.Time
}

func (User) TableName() string {
	return "user"
}
