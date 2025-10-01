package handler

import (
	"net/http"
	"strconv"

	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
)

func LogRetrievalHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		deploymentId := c.GetString("deploy_id")
		cursorStr := c.GetString("cursor")
		var cursor *int64
		if cursorStr != "" && cursorStr != "0" {
			if cursorVal, err := strconv.ParseInt(cursorStr, 10, 64); err == nil {
				cursor = &cursorVal
			}
		}
		res, err := srv.LogRetrievalService(c, deploymentId, cursor)
		if err != nil {
			if err.Error() == constants.UNAUTHORIZED_PROJECT_ACCESS {
				c.JSON(http.StatusUnauthorized, gin.H{"error": constants.UNAUTHORIZED_PROJECT_ACCESS})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
