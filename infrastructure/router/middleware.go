package router

import (
	"net/http"
	"store/domain"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	tokenUsecase domain.TokenUsecase
}

func NewMiddleware(tokenusecase domain.TokenUsecase) *Middleware {
	return &Middleware{
		tokenUsecase: tokenusecase,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := m.tokenUsecase.IsValidRequest(c.Request)
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
