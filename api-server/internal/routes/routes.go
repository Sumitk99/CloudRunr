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
	router.POST(
		"/deploy",
		middleware.ValidateDeployReq,
		handler.DeployReqHandler(srv))

}
