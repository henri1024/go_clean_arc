package router

import (
	"store/app/registry"
	"store/user/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(reg registry.Registry) *gin.Engine {
	r := gin.Default()

	UserEndpoint(r, reg.NewUserController())

	return r
}

func UserEndpoint(r *gin.Engine, handler controller.UserController) {
	r.POST("/api/users", handler.SignUp)
	r.POST("/api/users/login", handler.Login)
}
