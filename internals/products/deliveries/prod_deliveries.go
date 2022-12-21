package deliveries

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type productsController struct {
	ProductUse entities.ProductUsecase
}

func NewProductsDeliveries(r fiber.Router, productsUse entities.ProductUsecase) {
	deliveries := &productsController{
		ProductUse: productsUse,
	}
	r.Post("/create", deliveries.AddProduct)
	r.Get("/filters", deliveries.GetProduct)
}

func (h *productsController) AddProduct(c *fiber.Ctx) error {
	req := new(entities.ProductReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"status": fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message": err.Error(),
		})
	}

	res, err := h.ProductUse.AddProduct(req)
	if err != nil {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status": fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message": err.Error(),
			"result": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "OK",
			"status_code": fiber.StatusOK,
			"message": "",
			"result":res,
	})
}

func (h *productsController) GetProduct(c *fiber.Ctx) error {
	queryFilters := new(entities.ParamsFilters)

	if err := c.QueryParser(queryFilters); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":		"Internal Server Error",
			"status_code":	fiber.StatusInternalServerError,
			"message":		err.Error(),
			"result":		nil,
		})
	}

	res, err := h.ProductUse.GetProduct(queryFilters)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":      "Internal Server Error",
			"status_code": fiber.StatusInternalServerError,
			"message":     err.Error(),
			"result":      nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "OK",
		"status_code": fiber.StatusOK,
		"message":     nil,
		"result":      res,
	})

}