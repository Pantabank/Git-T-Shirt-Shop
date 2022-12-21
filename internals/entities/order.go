package entities

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrders(req *OrdersReq2)(*OrdersRes2, error)
	GetOrder(req *QueryParams)([]*OrdersRes, error)
}

type OrderUsecase interface {
	CreateOrders(req *OrdersReq2)(*OrdersRes2, error)
	GetOrder(req *QueryParams)([]*OrdersRes, error)
}

type QueryParams struct {
	Sdate 	string	`query:"sdate"`
	Edate 	string	`query:"edate"`
	Status	string	`query:"status"`
	Page 	int		`query:"page"`
	PerPage	int		`query:"per_page"`
}

type OrdersReq2 struct {
	Id			uuid.UUID `json:"id" db:"id"`
	Shipping 	*Shipping `json:"shipping_address"`
	Product 	ProductItem `json:"products"`
}

type OrdersRes2 struct {
	Id			int `json:"id" db:"id"`
	Shipping 	*Shipping `json:"shipping_address"`
	OrderID		int `json:"order_id" db:"order_id"`
	Product 	ProductItem `json:"products" db:"products"`
	Price 		int `json:"price" db:"price"`
	Qty 		int `json:"qty" db:"qty"`
}

type ProductItem struct {
	Item		[]Product `json:"item"`
}

type Shipping struct {
	FirstName 	string `json:"first_name" db:"first_name"`
	LastName	string `json:"last_name" db:"last_name"`
	SubDistrict	string `json:"sub_district" db:"sub_district"`
	District	string `json:"district" db:"district"`
	Province	string `json:"province"  db:"province" `
	Postcode	int `json:"postcode" db:"postcode"`
	Tel			string `json:"tel" db:"tel"`
}

type Product struct {
	Id	int `json:"id" db:"id"`
	Gender	string `json:"gender" db:"gender"`
	StyleType	string `json:"style_type" db:"style_type"`
	StyleDetail	string `json:"style_detail" db:"style_detail"`
	Size	string `json:"size" db:"size"`
	Price	int `json:"price" db:"price"`
	Qty	int `json:"qty" db:"qty"`
}

type OrdersReq struct {
	Id				uuid.UUID `json:"id" db:"id"`
	ProductID		int `json:"product_id" db:"product_id"`
	Gender			string `json:"gender" db:"gender"`
	StyleType		string `json:"style_type" db:"style_type"`
	StyleDetail		string `json:"style_detail" db:"style_detail"`
	Size			string `json:"size" db:"size"`
	Price			int `json:"price" db:"price" `
	ShippingAddress	AddressItem `json:"shipping_address" db:"shipping_address"`
	Status			string `json:"status" db:"status" sql:"type:statuses"`
	CreatedAt		time.Time `json:"created_at" db:"created_at"`
	Enable			bool `json:"enable" db:"enable"`
}

type OrdersRes struct {
	Id			int `json:"id" db:"id"`
	ProductID	int `json:"product_id" db:"product_id"`
	Gender		string `json:"gender" db:"gender"`
	StyleType	string `json:"style_type" db:"style_type"`
	StyleDetail	string `json:"style_detail" db:"style_detail"`
	Size		string `json:"size" db:"size"`
	Price		int `json:"price" db:"price" `
	ShippingAddress	AddressItem `json:"shipping_address" db:"shipping_address"`
	Status		string `json:"status" db:"status" sql:"type:statuses"`
	CreatedAt	time.Time `json:"created_at" db:"created_at"`
	Enable		bool `json:"enable" db:"enable"`
}

type AddressItem struct {
	First_name		string `json:"first_name"`
	Last_name		string `json:"last_name"`
	Sub_district	string `json:"sub_district"`
	District		string `json:"district"`
	Province		string `json:"province"`
	Postcode		int	`json:"postcode"`
	Tel				string `json:"tel"`
}



func (s ProductItem) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ProductItem) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s)
}