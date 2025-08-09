package handler

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"net/http"
)

func DeployReqHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form models.DeployReq
		err := c.BindJSON(&form)

		if form.ProjectID == nil || len(*form.ProjectID) == 0 {
			newSlug := slug.Make(*form.GitUrl)
			form.ProjectID = &newSlug
		}
		err = srv.ECSClient.SpinUpContainer(form.ProjectID, form.GitUrl, form.Framework, form.DistFolder)

		if err != nil {
			c.JSON(http.StatusInternalServerError, &models.DeployRes{
				Status: constants.STATUS_FAILED,
				Error:  err.Error(),
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, models.DeployRes{
			Status: constants.STATUS_QUEUED,
			Url:    *form.ProjectID,
		})
	}
}
