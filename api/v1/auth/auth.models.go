package v1_auth

import "github.com/adrieljss/golighter/models"

type userRegister struct {
	Username string `json:"username" validate:"required,max=30" faker:"username"`
	Email    string `json:"email" validate:"required,email" faker:"email"`
	Password string `json:"password" validate:"required" faker:"password"`
}

type userLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type userResponse struct {
	User  *models.User  `json:"user"`
	Token tokenResponse `json:"token"`
}

type userRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
