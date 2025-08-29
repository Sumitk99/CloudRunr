package handler

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeployReqHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		deployform, _ := c.Get("deploy_req")

		form := deployform.(models.DeployReq)
		res, err := srv.DeploymentService(c, &form.ProjectID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, &models.DeployRes{
				Status: constants.STATUS_FAILED,
				Error:  err.Error(),
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, models.DeployRes{
			Status:       constants.STATUS_QUEUED,
			Url:          form.ProjectID,
			DeploymentID: *res,
		})
	}
}

func NewProjectHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		projectForm, _ := c.Get("project_req")
		project := projectForm.(models.NewProjectReq)
		deploymentId, err := srv.NewProjectService(c, &project)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, models.NewProjectRes{
			DeploymentId: *deploymentId,
		})
	}
}

func GetUserProjectsHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := srv.GetUserProjectsService(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"projects": res})
	}
}
