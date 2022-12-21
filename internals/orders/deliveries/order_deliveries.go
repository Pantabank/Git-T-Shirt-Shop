package deliveries

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type ordersController struct {
	OrderUse entities.OrderUsecase
}

func NewOrdersDeliveries(r fiber.Router, ordersUse entities.OrderUsecase) {
	deliveries := &ordersController{
		OrderUse: ordersUse,
	}
	r.Post("", deliveries.CreateOrders)
	r.Get("/filters", deliveries.GetOrder)
}

func (h *ordersController) CreateOrders(c *fiber.Ctx) error {
	req := new(entities.OrdersReq2)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"status": fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message": err.Error(),
		})
	}

	if req.Shipping == nil {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"status": fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message": "non-shipping",
		})
	}

	if len(req.Product.Item) < 1 {
		return c.Status(fiber.ErrBadGateway.Code).JSON(fiber.Map{
			"status": fiber.ErrBadRequest.Message,
			"status_code": fiber.ErrBadRequest.Code,
			"message": "error product",
		})
	}



	res, err := h.OrderUse.CreateOrders(req)
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

func (h *ordersController) GetOrder(c *fiber.Ctx) error {
	queryParams := new(entities.QueryParams)

	if err := c.QueryParser(queryParams); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":		"Internal Server Error",
			"status_code":	fiber.StatusInternalServerError,
			"message":		err.Error(),
			"result":		nil,
		})
	}

	res, err := h.OrderUse.GetOrder(queryParams)
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