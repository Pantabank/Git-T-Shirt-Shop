package repositories

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type authRepo struct {
	Db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) entities.AuthRepository{
	return &authRepo{
		Db: db,
	}
}

func (r *authRepo) SignUsersAccessToken(req *entities.UsersData) (string, error) {
	claims := entities.UsersClaims{
		Id:						req.Id,
		Username: 				req.Username,
		Role: 					req.Role,
		RegisteredClaims: 		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: "access_token",
			Subject: "users_access_token",
			ID: uuid.NewString(),
			Audience: []string{"users"},
		},
	}

	mySigningKey := os.Getenv("JWT_SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return ss, nil
}

func (r *authRepo) FindOneUser(username string)(*entities.UsersData, error){
	query := `
		SELECT "id", "username", "password", "role" FROM "users" WHERE "username" = $1;
	`

	res := new(entities.UsersData)
	if err := r.Db.Get(res, query, username); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, user not found")
	}
	return res, nil
}

func (r *authRepo) Register(req *entities.RegisterReq, pass_gen string)(*entities.RegisterRes, error){
	query := `
		INSERT INTO "users"("username", "password", "role")
		VALUES ($1, $2, $3)
		RETURNING "id", "username", "role";
	`
	result := new(entities.RegisterRes)
	rows, err := r.Db.Queryx(query, req.Username, pass_gen, req.Role)

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		for rows.Next() {
			if err := rows.StructScan(result); err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
		}

		return result, nil
}