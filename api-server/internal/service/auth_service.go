package service

import (
	"errors"
	"github.com/Sumitk99/CloudRunr/api-server/internal/constants"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"log"
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
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24*30)).Unix(),
		},
	}

	refreshClaims := &models.SignedDetails{ // used to get a new token if a token expires
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	singedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SECRET_KEY)

	singedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(SECRET_KEY)
	if err != nil {
		return
	}
	return
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func (srv *Service) SignUpService(req *models.SignUpReq) (*models.User, error) {
	password := HashPassword(req.Password)
	NewUser := &models.User{
		UserID:   ksuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	if err := srv.Repo.SignUpRepository(NewUser); err != nil {
		return nil, err
	}

	return NewUser, nil
}

func (srv *Service) LoginService(email, password *string) (*models.LoginResponse, error) {
	user, err := srv.Repo.GetUserByMail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password))
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("email or passsword is incorrect")
	}

	token, refreshtoken, err := srv.GenerateAllTokens(user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		UserID:       &user.UserID,
		Name:         &user.Name,
		Email:        &user.Email,
		GithubID:     &user.GithubID,
		Token:        &token,
		RefreshToken: &refreshtoken,
	}, nil
}
