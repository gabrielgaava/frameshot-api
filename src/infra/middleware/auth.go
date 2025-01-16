package middleware

import (
	"example/web-service-gin/src/utils"
	"github.com/gin-gonic/gin"
)

func JwtServiceMiddleware(jwksURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtService, err := utils.NewJwtService(jwksURL)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "failed to initialize JwtService"})
			return
		}

		// Adds JwtService on Gin context
		c.Set("jwtService", jwtService)
		c.Next()
	}
}
