package middleware

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var validate = validator.New()

func ValidateDeployReq(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectId := c.Param("project_id")
		if len(projectId) == 0 {
			c.JSON(http.StatusBadRequest,
				&models.DeployRes{
					Status: constants.STATUS_FAILED,
					Error:  "No projectId provided",
				},
			)
			c.Abort()
			return
		}

		c.Set("deploy_req", models.DeployReq{ProjectID: projectId})
	}

}

func ValidateNewProjectReq(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var form models.NewProjectReq
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
				c.Set("project_req", form)
				c.Next()
				return
			}
		}

		c.Abort()
		return

	}
}

func ValidateSignUpReq(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.SignUpReq
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json Req"})
			c.Abort()
			return
		}
		log.Println(user)

		// Validate the struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			log.Println(validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
			c.Abort()
			return
		}

		exists, err := srv.Repo.CheckUserExists(&user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking email availability"})
			c.Abort()
			return
		}

		if *exists {
			c.JSON(http.StatusConflict, gin.H{"error": "There's already an account registered with this email. Please log in."})
			c.Abort()
			return
		}
		c.Set("validated_signup_req", user)

		log.Println("SignUp Req Validated")
		c.Next()
	}
}

func ValidateLoginReq(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials models.LoginReq
		err := c.BindJSON(&credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing request"})
			c.Abort()
			return
		}
		validationErr := validate.Struct(credentials)
		if validationErr != nil {
			log.Println(validationErr)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
			c.Abort()
			return
		}

		c.Set("login_req", credentials)
		c.Next()
	}
}

func ValidateLogRetrievalReq(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		deploymentId := c.Param("deploy_id")
		if deploymentId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide deployment ID"})
			c.Abort()
			return
		}

		offset := c.Param("offset")
		if offset == "" {
			offset = "1"
		}
		c.Set("deploy_id", deploymentId)
		c.Set("offset", offset)
		c.Next()
	}
}
