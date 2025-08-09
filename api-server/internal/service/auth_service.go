package service

import (
	"errors"
	"fmt"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var SECRET_KEY = []byte("depsecret")

func ValidateToken(signedToken string) (claims *models.SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&models.SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*models.SignedDetails)
	if !ok {
		return nil, errors.New(constants.INVALID_TOKEN)
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New(constants.TOKEN_EXPIRED)
	}
	return claims, nil
}

func (srv *Service) GenerateAllTokens(user *models.User) (singedToken string, singedRefreshToken string, err error) {
	claims := &models.SignedDetails{
		Email:    user.Email,
		Name:     user.Name,
		UserID:   user.UserID,
		GithubID: user.GithubID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	fmt.Println("created claims")

	refreshClaims := &models.SignedDetails{ // used to get a new token if a token expires
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	fmt.Println("created refresh claims")

	singedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)
	//singedToken, err = jwt.NewWithClaims(jwt.SigningMethodNone, claims).SignedString(jwt.UnsafeAllowNoneSignatureType)
	fmt.Println("created tokens")
	fmt.Println(err)
	//singedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodNone, refreshClaims).SignedString(jwt.UnsafeAllowNoneSignatureType)

	singedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	fmt.Println("created refresh tokens")
	fmt.Println(err)
	if err != nil {
		return
	}
	fmt.Println("over")
	return
}
