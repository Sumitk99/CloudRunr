package middleware

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func ValidateDeployReq(c *gin.Context) {
	var form models.DeployReq
	err := c.BindJSON(&form)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Parsing Request"})
		c.Abort()
		return
	}

	if form.GitUrl == nil ||
		(!strings.HasPrefix(*form.GitUrl, constants.GITHUB_URL_PREFIX_1) &&
			!strings.HasPrefix(*form.GitUrl, constants.GITHUB_URL_PREFIX_2)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.INVALID_GITHUB_URL_MESSAGE})
		c.Abort()
		return
	}

	for _, valid := range constants.VALID_FRAMEWORKS {
		if *form.Framework == valid {
			c.Next()
			return
		}
	}

	c.Abort()
	return
}
