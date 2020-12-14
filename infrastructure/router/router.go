package router

import (
	"go_clean_arc/app"

	"github.com/gin-gonic/gin"
)

func NewRouter(app *app.App) *gin.Engine {
	r := gin.Default()

	r.POST("/api/user/new", app.UserController.SignUp)
	r.POST("/api/user/login", app.UserController.SignIn)
	r.GET("/api/user/profile", app.UserController.AuthMiddleware(), app.UserController.GetProfile)

	return r
}
