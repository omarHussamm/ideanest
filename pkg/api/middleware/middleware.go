package middleware

import (
	"demo/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccessJwtAuthMiddleware(c *gin.Context) {
	err := utils.AccessTokenValid(c)
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()

}
