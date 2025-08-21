package routes

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/handler"
	"github.com/Sumitk99/CloudRunr/api-server/internal/middleware"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, srv *service.Service) {
	router.POST(
		"/signup",
		middleware.ValidateSignUpReq(srv),
		handler.SignUpHandler(srv),
	)
	router.POST(
		"/login",
		middleware.ValidateLoginReq(srv),
		handler.LoginHandler(srv))
	router.GET(
		"/user",
		middleware.Authenticate(),
		handler.GetUser(srv),
	)
	router.POST(
		"/project",
		middleware.Authenticate(),
		middleware.ValidateNewProjectReq(srv),
		handler.NewProjectHandler(srv),
	)
	router.POST(
		"/deploy/:project_id",
		middleware.Authenticate(),
		middleware.ValidateDeployReq(srv),
		handler.DeployReqHandler(srv))
}
