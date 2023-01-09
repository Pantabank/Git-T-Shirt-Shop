package entities

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pantabank/t-shirts-shop/configs"

	"github.com/google/uuid"
)

type AuthRepository interface {
	SignUsersAccessToken(req *UsersData)(string, error)
	Register(req *RegisterReq, pass_gen string)(*RegisterRes, error)
	FindOneUser(username string)(*UsersData, error)
}

type AuthUsecase interface {
	Login(cfg *configs.Configs, req *UsersCredentials)(*UsersLoginRes, error)
	Register(req *RegisterReq)(*RegisterRes, error)
}

type UsersCredentials struct {
	Username string `json:"username" db:"username" form:"username"`
	Password string `json:"password" db:"password" form:"password"`
}

type UsersData struct {
	Id 			int `json:"id" db:"id"`
	Username 	string `json:"username" db:"username"`
	Password 	string `json:"password" db:"password"`
	Role		string `json:"role" db:"role"`
}

type UsersClaims struct {
	Id 			int `json:"id" db:"id"`
	Username 	string `json:"username" db:"username"`
	Role		string `json:"role" db:"role"`
	jwt.RegisteredClaims
}

type UsersLoginRes struct {
	AccessToken	string `json:"access_token"`
}

type RegisterReq struct {
	Id			uuid.UUID `json:"id" db:"id"`
	Username	string `json:"username" db:"username"`
	Password	string `json:"password" db:"password"`
	Role		string `json:"role" db:"role" sql:"type:roles"`
}

type RegisterRes struct {
	Id			string `json:"id" db:"id"`
	Username	string `json:"username" db:"username"`
	Role		string `json:"role" db:"role"`
}