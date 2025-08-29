package middleware

import (
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		log.Println("Token : ", clientToken)
		if clientToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf(constants.NO_TOKEN)})
			c.Abort()
			return
		}

		claims, err := service.ValidateToken(clientToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("user_id", claims.UserID)
		c.Set("github_id", claims.GithubID)
		c.Next()
	}
}
