package handler

import (
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func SignUpHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("calling signup service")
		validatedUser, _ := c.Get("validated_signup_req")

		var user models.SignUpReq
		user = validatedUser.(models.SignUpReq)
		NewUser, err := srv.SignUpService(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		token, refreshToken, err := srv.GenerateAllTokens(NewUser)
		res := &models.SignUpResponse{
			UserID:       &NewUser.UserID,
			Name:         &NewUser.Name,
			Email:        &NewUser.Email,
			Token:        &token,
			RefreshToken: &refreshToken,
		}

		c.JSON(http.StatusCreated, res)
	}
}

func LoginHandler(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, _ := c.Get("login_req")
		loginReq := req.(models.LoginReq)
		user, err := srv.LoginService(loginReq.Email, loginReq.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(http.StatusAccepted, user)
	}
}

func GetUser(srv *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := c.GetString("email")
		user, err := srv.Repo.GetUserByMail(&userEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Error getting user data %s\n", err.Error()),
			})
			c.Abort()
			return
		}
		user.Password = ""
		c.JSON(http.StatusOK, user)
	}
}
