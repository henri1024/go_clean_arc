package router

import (
	"userauth/app/registry"
	"userauth/user/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(reg registry.Registry) *gin.Engine {
	r := gin.Default()

	middleware := NewMiddleware(reg.NewUserController().TokenUsecase)

	UserEndpoint(r, reg.NewUserController(), middleware)

	return r
}

func UserEndpoint(r *gin.Engine, handler controller.UserController, middleware *Middleware) {
	r.POST("/api/users/signup", handler.SignUp)
	r.POST("/api/users/login", handler.Login)
	r.POST("/api/users/logout", middleware.AuthMiddleware(), handler.LogOut)
	r.POST("/api/users/refresh", handler.Refresh)
	r.PUT("/api/users/updatepassword", middleware.AuthMiddleware(), handler.ChangePassword)
}
