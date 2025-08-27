package handler

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func LogRetrievalHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		deploymentId := c.GetString("deploy_id")
		offsetStr := c.GetString("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 1
		}
		res, err := srv.LogRetrievalService(c, deploymentId, offset)
		if err != nil {
			if err.Error() == constants.UNAUTHORIZED_PROJECT_ACCESS {
				c.JSON(http.StatusUnauthorized, gin.H{"error": constants.UNAUTHORIZED_PROJECT_ACCESS})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, models.LogRetrievalResponse{Data: res})
	}
}
