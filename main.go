package main

import (
	"blog-go/src/config"
	"blog-go/src/controller"
	"github.com/kataras/iris/v12"
	"github.com/rs/cors"
)

func main() {
	app := newApp()
	app.Listen(":8081")
}

func newApp() *iris.Application {
	app := iris.Default()

	//+--------------------------------+ CORS config
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	app.WrapRouter(c.ServeHTTP)
	//+--------------------------------+
	app.Post("/login", controller.Login)

	authApi := app.Party("")
	authApi.Use(config.J.Verify)
	authApi.Get("/currentUserName", controller.CurrentUserName)
	authApi.Get("/currentUserId", controller.CurrentUserId)
	authApi.Get("/logout", controller.Logout)

	return app
}

func index(ctx iris.Context) {
	ctx.HTML("<h1>Welcome</h1>")
}
