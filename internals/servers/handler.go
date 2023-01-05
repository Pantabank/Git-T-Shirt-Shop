package servers

import (
	_productsHttp "github.com/pantabank/t-shirts-shop/internals/products/deliveries"
	_productsRepository "github.com/pantabank/t-shirts-shop/internals/products/repositories"
	_productsUsecase "github.com/pantabank/t-shirts-shop/internals/products/usecases"

	_ordersHttp "github.com/pantabank/t-shirts-shop/internals/orders/deliveries"
	_ordersRepository "github.com/pantabank/t-shirts-shop/internals/orders/repositories"
	_ordersUsecase "github.com/pantabank/t-shirts-shop/internals/orders/usecases"

	_authHttp "github.com/pantabank/t-shirts-shop/internals/auth/deliveries"
	_authRepository "github.com/pantabank/t-shirts-shop/internals/auth/repositories"
	_authUsecase "github.com/pantabank/t-shirts-shop/internals/auth/usecases"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	v1 := s.App.Group("/v1")

	productsGroup := v1.Group("/products")
	productsRepository := _productsRepository.NewProductsRepository(s.Db)
	productsUsecase := _productsUsecase.NewProductUsecase(productsRepository)
	_productsHttp.NewProductsDeliveries(productsGroup, productsUsecase)

	ordersGroup := v1.Group("/orders")
	ordersRepository := _ordersRepository.NewOrdersRepository(s.Db)
	ordersUsecase := _ordersUsecase.NewOrderUsecase(ordersRepository)
	_ordersHttp.NewOrdersDeliveries(ordersGroup, ordersUsecase)

	authGroup := v1.Group("/auth")
	authRepository := _authRepository.NewAuthRepository(s.Db)
	authUsecase := _authUsecase.NewAuthUsecase(authRepository)
	_authHttp.NewAuthController(authGroup, s.Cfg, authUsecase)

	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status": fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message": "error, end point not found",
			"result": nil,
		})
	})
	return nil
}