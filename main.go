package main

import (
	"blog-go/src/controller"
	"blog-go/src/domain"
	"fmt"
	"log"

	"github.com/kataras/iris/v12"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:88888888@tcp(127.0.0.1:3306)/vueblog2?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("db conn fail.")
	} else {
		log.Println("db conn success.")
	}
	user := domain.User{}
	db.Preload("Roles").First(&user, 6)
	//db.Debug().Preload("Roles").First(&user, 6)
	fmt.Println(user.Roles)
	app := newApp()

	app.Listen(":8081")
}

func newApp() *iris.Application {
	app := iris.Default()
	app.Get("/currentUserName", controller.CurrentUserName)
	app.Get("/currentUserId", controller.CurrentUserId)
	app.Post("/login", controller.Login)
	//app.HandleDir("/", iris.Dir("./static"))
	return app
}

func index(ctx iris.Context) {
	ctx.HTML("<h1>Index Page</h1>")
}
