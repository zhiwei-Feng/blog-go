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
	app.Logger().SetLevel("DEBUG")

	//+--------------------------------+ CORS config
	c := cors.AllowAll()
	app.WrapRouter(c.ServeHTTP)
	//+--------------------------------+
	app.Post("/login", controller.Login)

	authApi := app.Party("")
	authApi.Use(config.J.Verify)
	authApi.Get("/currentUserName", controller.CurrentUserName)
	authApi.Get("/currentUserId", controller.CurrentUserId)
	authApi.Get("/currentUserEmail", controller.CurrentUserEmail)
	authApi.Get("/isAdmin", controller.IsAdmin)
	authApi.Put("/updateUserEmail", controller.UpdateUserEmail)
	authApi.Get("/logout", controller.Logout)

	admin := authApi.Party("/admin")
	{
		admin.Get("/user", controller.GetUserByNickname)
		admin.Get("/user/{id:int}", controller.GetUserById)
		admin.Get("/roles", controller.GetAllRoles)
		admin.Put("/user/role", controller.UpdateUserRoles)
		admin.Put("/user/enabled", controller.SwitchUserEnableStatus)
		admin.Delete("/user/{id:int}", controller.DeleteUserById)

		cate := admin.Party("/category")
		{
			cate.Get("/all", controller.GetAllCategories)
			cate.Post("/", controller.AddNewCategory)
			cate.Put("/", controller.UpdateCategory)
			cate.Delete("/{ids:string}", controller.DeleteCateById)
		}

		arti := admin.Party("/article")
		{
			arti.Get("/all", controller.GetArticleByStateForAdmin)
		}
	}

	article := authApi.Party("/article")
	{
		article.Get("/dataStatistics", controller.GetDataStatistics)
		article.Get("/all", controller.GetArticleByState)
		article.Get("/{aid:int}", controller.GetArticleById)
		article.Post("/", controller.AddNewArticle)
	}

	return app
}
