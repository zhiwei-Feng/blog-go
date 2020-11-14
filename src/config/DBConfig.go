package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func init() {
	dsn := "root:88888888@tcp(127.0.0.1:3306)/vueblog2?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("db conn fail.")
	} else {
		log.Println("db conn success.")
		DB = db
	}
}
