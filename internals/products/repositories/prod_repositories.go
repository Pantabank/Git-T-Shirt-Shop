package repositories

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type productsRepo struct {
	Db *sqlx.DB
}

func NewProductsRepository(db *sqlx.DB) entities.ProductRepository {
	return &productsRepo{
		Db: db,
	}
}


func (r *productsRepo) AddProduct(req *entities.ProductReq) (*entities.ProductRes, error){
	query := `
		INSERT INTO "products"(
			"gender",
			"style_type",
			"style_detail",
			"size",
			"price",
			"enable"
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "id", "gender", "style_type", "style_detail", "size", "price", "enable";
	`

	user := new(entities.ProductRes)
	rows, err := r.Db.Queryx(query, strings.ToLower(req.Gender), strings.ToLower(req.StyleType), req.StyleDetail, strings.ToLower(req.Size), req.Price, `TRUE`)
	if err != nil {
		defer rows.Close()
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(user); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}

	return user, nil
}

func (r *productsRepo) GetProduct(req *entities.ParamsFilters)(list []*entities.ProductRes, err error){
	lists := make([]*entities.ProductRes, 0)
	query := `SELECT * FROM products WHERE enable=true`

	if req.Gender != "" {
		query += fmt.Sprintf(" AND gender='%v'", strings.ToLower(req.Gender))
	}

	if req.StyleType != "" {
		query += fmt.Sprintf(" AND style_type='%v'", strings.ToLower(req.StyleType))
	}

	if req.Size != "" {
		query += fmt.Sprintf(" AND size='%v'", strings.ToLower(req.Size))
	}

	pages := req.PerPage * (req.Page - 1)
	query += fmt.Sprintf(` LIMIT %d OFFSET %d`, req.PerPage, pages)

	rows, err := r.Db.Query(query); 
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		products := new(entities.ProductRes)
		err := rows.Scan(&products.Id, &products.Gender, &products.StyleType, &products.StyleDetail, &products.Size, &products.Price, &products.Enable)
		if err != nil{
			return nil, err
		}
		lists = append(lists, products)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return lists, nil
}