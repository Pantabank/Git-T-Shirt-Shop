package usecases

import (
	"errors"

	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type productsUse struct {
	ProductsRepo entities.ProductRepository
}

func NewProductUsecase(productsRepo entities.ProductRepository) entities.ProductUsecase {
	return &productsUse{
		ProductsRepo: productsRepo,
	}
}

func (u *productsUse) AddProduct(req *entities.ProductReq)(*entities.ProductRes, error) {
	user, err := u.ProductsRepo.AddProduct(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *productsUse) GetProduct(req *entities.ParamsFilters)([]*entities.ProductRes, error){
	info, err := p.ProductsRepo.GetProduct(req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return info, nil
}