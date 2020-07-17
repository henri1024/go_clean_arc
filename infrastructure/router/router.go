package router

import (
	"store/app/registry"
	"store/user/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(reg registry.Registry) *gin.Engine {
	r := gin.Default()

	// middleware := NewMiddleware(reg.NewUserController().TokenUsecase)

	UserEndpoint(r, reg.NewUserController())

	return r
}

func UserEndpoint(r *gin.Engine, handler controller.UserController) {
	r.POST("/api/users/signup", handler.SignUp)
	r.POST("/api/users/login", handler.Login)
	r.POST("/api/users/logout", handler.LogOut)
}
