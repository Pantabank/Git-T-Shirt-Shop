package entities

import (
	"github.com/google/uuid"
)

type ProductRepository interface {
	AddProduct(req *ProductReq) (*ProductRes, error)
	GetProduct(req *ParamsFilters)([]*ProductRes, error)
}

type ProductUsecase interface {
	AddProduct(req *ProductReq) (*ProductRes, error)
	GetProduct(req *ParamsFilters)([]*ProductRes, error)
}

type ParamsFilters struct {
	Gender 	string	`query:"gender"`
	StyleType 	string	`query:"style_type"`
	Size	string	`query:"size"`
	Page 	int		`query:"page"`
	PerPage	int		`query:"per_page"`
}

type ProductReq struct {
		Id			uuid.UUID `json:"id" db:"id"`
		Gender		string `json:"gender" db:"gender" sql:"type:genders"`
		StyleType	string `json:"style_type" db:"style_type"`
		StyleDetail	string `json:"style_detail" db:"style_detail"`
		Size		string `json:"size" db:"size"`
		Price		int `json:"price" db:"price" `
		Enable		bool `json:"enable" db:"enable"`
}

type ProductRes struct {
	Id			int `json:"id" db:"id"`
	Gender		string `json:"gender" db:"gender" sql:"type:genders"`
	StyleType	string `json:"style_type" db:"style_type"`
	StyleDetail	string `json:"style_detail" db:"style_detail"`
	Size		string `json:"size" db:"size"`
	Price		int `json:"price" db:"price"`
	Enable		bool `json:"enable" db:"enable"`
}

