package main

import (
	"blog-go/src/controller"
	"github.com/kataras/iris/v12"
	"github.com/rs/cors"
)

func main() {

	//user := domain.User{}
	//db.Preload("Roles").First(&user, 6)
	//db.Debug().Preload("Roles").First(&user, 6)
	//fmt.Println(user.Roles)
	app := newApp()

	app.Listen(":8081")
}

func newApp() *iris.Application {
	app := iris.Default()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	app.WrapRouter(c.ServeHTTP)

	app.Get("/currentUserName", controller.CurrentUserName)
	app.Get("/currentUserId", controller.CurrentUserId)
	app.Get("/logout", controller.Logout)
	app.Post("/login", controller.Login)
	//app.Get("/", index)
	//app.HandleDir("/", iris.Dir("./static"))
	return app
}

func index(ctx iris.Context) {
	ctx.HTML("<h1>Welcome</h1>")
}
