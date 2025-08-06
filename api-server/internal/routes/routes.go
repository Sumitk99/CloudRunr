package routes

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/handler"
	"github.com/Sumitk99/CloudRunr/api-server/internal/middleware"
	"github.com/Sumitk99/CloudRunr/api-server/internal/server"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, ecsConfig *server.ECSClusterConfig) {
	router.POST(
		"/deploy",
		middleware.ValidateDeployReq,
		handler.DeployReqHandler(ecsConfig))
}
