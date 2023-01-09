package usecases

import (
	"errors"
	"fmt"
	"strings"
	"time"

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

func (u *ordersUse) CreateOrders(req *entities.OrdersReq2, uid any)(*entities.OrdersRes2, error) {

	shipping, err := u.GetOrderID(req.Shipping, uid)
	if err != nil {
		return nil, err
	}

	times := time.Now()
	//order := new(entities.OrdersRes2)
	product := []entities.Product{}
	collections := make(map[string][]entities.Product)
	var totalQty, totalPrice int

	for _, v := range req.Product.Item {
		productRes, err := u.QueryCart(v.Id) 
		p := entities.Product{Id: productRes.Id, Gender: strings.ToLower(productRes.Gender), StyleType: strings.ToLower(productRes.StyleType), StyleDetail: productRes.StyleDetail, Size: strings.ToLower(productRes.Size), Price: productRes.Price, Qty: v.Qty, TotalPrice: productRes.Price * float64(v.Qty)}
		if err != nil {
			fmt.Println(err.Error())
		}
		totalQty += v.Qty
		totalPrice += int(productRes.Price) * v.Qty
		product = append(product, p)

		u.OrdersRepo.AddOrders(shipping.Id, productRes.Id, &p)
	}

	collections["item"] = append(collections["item"], product...)
	u.OrdersRepo.UpdateOrders(totalQty, totalPrice, shipping.Id)
	pd := entities.OrdersRes2{
		OrderID: 	shipping.Id,
		CustomerID: uid,
		Qty: 		totalQty,
		Price: 		totalPrice,
		Shipping: 	shipping.Shipping,
		Product: 	entities.ProductItem{product},
		Status: 	"placed_order",
		CreatedAt: 	times,
	}
	return &pd, nil
}

func (u *ordersUse) GetOrderID(req *entities.Shipping, uid any)(*entities.AddressesRes, error) {
	order, err := u.OrdersRepo.GetOrderID(req, uid)
	if err != nil {
		return nil, err
	}
	return order, nil 
}

func (u *ordersUse) QueryCart(id int)(*entities.Product, error) {
	order, err := u.OrdersRepo.QueryCart(id)
	if err != nil {
		return nil, err
	}
	return order, nil 
}

func (u *ordersUse) UpdateOrders(qty, price, order_id int) error {
	order := u.OrdersRepo.UpdateOrders(qty, price, order_id)
	return order
}

func (u *ordersUse) AddOrders(order_id, product_id int, product *entities.Product) error {
	order := u.OrdersRepo.AddOrders(order_id, product_id, product)
	return order
}

func (o *ordersUse) GetOrder(req *entities.QueryParams)([]*entities.GetOrderRes, error){
	info, err := o.OrdersRepo.GetOrder(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return info, nil
}