package usecases

import (
	"errors"

	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type ordersUse struct {
	OrdersRepo entities.OrderRepository
}

func NewOrderUsecase(ordersRepo entities.OrderRepository) entities.OrderUsecase {
	return &ordersUse{
		OrdersRepo: ordersRepo,
	}
}

func (u *ordersUse) CreateOrders(req *entities.OrdersReq2)(*entities.OrdersRes2, error) {
	order, err := u.OrdersRepo.CreateOrders(req)
	if err != nil {
		return nil, err
	}
	return order, nil 
}

func (o *ordersUse) GetOrder(req *entities.QueryParams)([]*entities.GetOrderRes, error){
	info, err := o.OrdersRepo.GetOrder(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return info, nil
}