package handler

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

//func GetUser() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		userId := c.Param("user_id")
//
//		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
//
//		var user models.User
//		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user) // using decode to make it string as golang doesn't understand json
//		defer cancel()
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusOK, user)
//	}
//}
