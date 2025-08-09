package handler

import (
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(userPassword string, providedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}

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

//func Login() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
//		defer cancel()
//		var user models.User
//		var foundUser models.User //to check if the person already exists in the database
//		if err := c.BindJSON(&user); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			return
//		}
//		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"err": "User with the Provided Credentials does not exist"})
//			return
//		}
//		fmt.Println(err)
//		defer cancel()
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
//			return
//		}
//		fmt.Println(" user not found")
//		if foundUser.Email == nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
//			return
//		}
//		fmt.Println("verifying password")
//		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
//		defer cancel()
//		fmt.Println("verified password")
//		if passwordIsValid != true {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
//			return
//		}
//
//		token, refreshToken, _ := helper.GenerateAllTokens(
//			*foundUser.Email,
//			*foundUser.FirstName,
//			*foundUser.LastName,
//			*foundUser.UserType,
//			*&foundUser.UserId)
//
//		helper.UpdateAllTokens(token, refreshToken, foundUser.UserId)
//
//		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserId}).Decode(&foundUser)
//
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//			return
//		}
//		c.JSON(http.StatusOK, foundUser)
//	}
//}

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
