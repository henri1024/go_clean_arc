package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uc *userController) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := uc.authUsecase.IsValidRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": http.StatusUnauthorized,
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
