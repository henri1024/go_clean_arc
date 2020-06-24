package router

import (
	"clean_arc/interface/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter(controller controllers.AppController) *gin.Engine {
	r := gin.Default()

	r.POST("/api/users", controller.SaveUser)

	return r
}
