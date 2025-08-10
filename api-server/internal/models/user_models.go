package models

import "github.com/dgrijalva/jwt-go"

type SignUpReq struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignUpResponse struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	UserID       *string `json:"user_id"`
	Token        *string `json:"token"`
	RefreshToken *string `json:"refresh_token"`
}

type LoginReq struct {
	Email    *string `json:"email" validate:"email,required"`
	Password *string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Name         *string `json:"name"`
	Email        *string `json:"email"`
	UserID       *string `json:"user_id"`
	Token        *string `json:"token"`
	RefreshToken *string `json:"refresh_token"`
	GithubID     *string `json:"github_id"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserID   string `json:"user_id"`
	Password string `json:"password,omitempty"`
	GithubID string `json:"github_id"`
}

type SignedDetails struct {
	Email    string
	Name     string
	UserID   string
	GithubID string
	jwt.StandardClaims
}
