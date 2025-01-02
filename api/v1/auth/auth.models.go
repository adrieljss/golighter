package v1_auth

import "github.com/adrieljss/golighter/models"

type UserRegister struct {
	Username string `json:"username" validate:"required,max=30,username" fake:"{username}"`
	Email    string `json:"email" validate:"required,email" fake:"{email}"`
	Password string `json:"password" validate:"required" fake:"{password}"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserResponse struct {
	User  *models.User  `json:"user"`
	Token TokenResponse `json:"token"`
}

type UserRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
