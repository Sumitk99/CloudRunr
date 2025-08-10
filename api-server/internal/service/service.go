package service

import (
	"errors"
	"github.com/Sumitk99/CloudRunr/api-server/internal/models"
	"github.com/Sumitk99/CloudRunr/api-server/internal/repository"
	"github.com/Sumitk99/CloudRunr/api-server/internal/server"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Service struct {
	Repo      *repository.Repository
	ECSClient *server.ECSClusterConfig
}

func NewService(repo *repository.Repository, ecsConfig *server.ECSClusterConfig) *Service {
	return &Service{
		Repo:      repo,
		ECSClient: ecsConfig,
	}
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
	log.Println("provided : ", password)
	log.Println("crypted : ", user.Password)
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
