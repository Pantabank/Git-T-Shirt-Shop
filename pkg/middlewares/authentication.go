package middlewares

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JwtAuthentication(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if accessToken == "" {
			log.Println("error, authorization header is empty.")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": 		"unauthorized",
				"status_code":	fiber.StatusUnauthorized,
				"message":		"error, unauthorized",
				"rresult":		nil,	
			})
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error, unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
				"status": 		fiber.ErrUnauthorized.Message,
				"status_code":	fiber.ErrUnauthorized.Code,
				"message":		err.Error(),
				"result":		nil,
			})
		}

		if token.Claims.(jwt.MapClaims)["role"] != role {
			return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
				"status": 		fiber.ErrUnauthorized.Message,
				"status_code":	fiber.ErrUnauthorized.Code,
				"message":		"your role cannot access this endpoint.",
				"result":		nil,
			})
		}

		if err != nil {
			log.Println(err.Error())
			return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
				"status": 		fiber.ErrUnauthorized.Message,
				"status_code":	fiber.ErrUnauthorized.Code,
				"message":		err.Error(),
				"result":		nil,
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims)
			c.Locals("user_id", claims["user_id"])
			c.Locals("username", claims["username"])
			return c.Next()
		}

		return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"status": 		fiber.ErrUnauthorized.Message,
			"status_code":	fiber.ErrUnauthorized.Code,
			"message":		"error, unauthorized",
			"result":		nil,
		})
	}
}

func UserData(c *fiber.Ctx)error {
	accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		

		secretKey := os.Getenv("JWT_SECRET_KEY")
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error, unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		a := token.Claims.(jwt.MapClaims)["id"]
		fmt.Println(a)
		return err
}