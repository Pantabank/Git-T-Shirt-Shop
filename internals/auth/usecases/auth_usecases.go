package usecases

import (
	"errors"
	"fmt"

	"github.com/pantabank/t-shirts-shop/configs"
	"github.com/pantabank/t-shirts-shop/internals/entities"
	"golang.org/x/crypto/bcrypt"
)

type authUse struct {
	AuthRepo		entities.AuthRepository
}

func NewAuthUsecase(authRepo entities.AuthRepository)entities.AuthUsecase {
	return &authUse{
		AuthRepo:	authRepo,
	}
}

func (u *authUse) Login(cfg *configs.Configs, req *entities.UsersCredentials)(*entities.UsersLoginRes, error){
	user, err := u.AuthRepo.FindOneUser(req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, password is invalid")
	}

	token, err := u.AuthRepo.SignUsersAccessToken(user)
	if err != nil {
		return nil, err
	}
	res := &entities.UsersLoginRes{
		AccessToken: token,
	}
	return res, nil
}