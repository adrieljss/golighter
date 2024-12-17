package v1_users

import (
	"github.com/adrieljss/golighter/platform"
)

type UserUsecase struct {
	UserRepo *userRepository
	*platform.Application
}

func NewUserUsecase(app *platform.Application) *UserUsecase {
	return &UserUsecase{
		UserRepo:    NewUserRepository(app),
		Application: app,
	}
}
