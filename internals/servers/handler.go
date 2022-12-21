package servers

import (
	_productsHttp "github.com/pantabank/t-shirts-shop/internals/products/deliveries"
	_productsRepository "github.com/pantabank/t-shirts-shop/internals/products/repositories"
	_productsUsecase "github.com/pantabank/t-shirts-shop/internals/products/usecases"

	_ordersHttp "github.com/pantabank/t-shirts-shop/internals/orders/deliveries"
	_ordersRepository "github.com/pantabank/t-shirts-shop/internals/orders/repositories"
	_ordersUsecase "github.com/pantabank/t-shirts-shop/internals/orders/usecases"

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