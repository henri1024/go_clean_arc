package router

import (
	"clean_arc/interface/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(controller controller.AppController) *gin.Engine {
	r := gin.Default()

	r.POST("/api/users", controller.SaveUser)
	r.POST("/api/users/login", controller.GetUserByEmailAndPassowrd)

	return r
}
