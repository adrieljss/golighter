package v1_auth

import (
	"context"
	"errors"

	v1_users "github.com/adrieljss/golighter/api/v1/users"
	"github.com/adrieljss/golighter/models"
	"github.com/adrieljss/golighter/platform"
	"github.com/adrieljss/golighter/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	UserUsecase *v1_users.UserUsecase
	*platform.Application
}

func NewAuthUsecase(app *platform.Application) *AuthUsecase {
	return &AuthUsecase{
		UserUsecase: v1_users.NewUserUsecase(app),
		Application: app,
	}
}

func (u *AuthUsecase) GenerateTokenPair(user *models.User) (string, string, error) {
	userTokenClaim := utils.UserTokenClaim{
		UID:  user.UID,
		Role: user.Role,
	}
	accessToken, err := utils.GenerateJWT(&userTokenClaim, u.Env.JWTSecretAccessToken, u.Env.JWTAccessTokenTTL)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateJWT(&userTokenClaim, u.Env.JWTSecretRefreshToken, u.Env.JWTRefreshTokenTTL)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func (u *AuthUsecase) CreateUser(ctx context.Context, user *userRegister) (*userResponse, error) {
	newUser, err := u.UserUsecase.UserRepo.Create(ctx, user.Username, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := u.GenerateTokenPair(newUser)
	if err != nil {
		return nil, err
	}

	return &userResponse{
		User: newUser,
		Token: tokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (u *AuthUsecase) LoginUser(ctx context.Context, user *userLogin) (*userResponse, error) {
	newUser, err := u.UserUsecase.UserRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(newUser.PasswordHash), []byte(user.Password)); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := u.GenerateTokenPair(newUser)
	if err != nil {
		return nil, err
	}

	return &userResponse{
		User: newUser,
		Token: tokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// RefreshAccessToken refreshes the access token, returns the new access token
func (u *AuthUsecase) RefreshAccessToken(ctx context.Context, refresh_token string) (string, error) {
	claims, err := utils.ValidateJWT(refresh_token, u.Env.JWTSecretAccessToken)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateJWT(&claims.UserTokenClaim, u.Env.JWTSecretAccessToken, u.Env.JWTAccessTokenTTL)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
