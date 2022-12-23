package repositories

import (
	//"encoding/json"

	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pantabank/t-shirts-shop/internals/entities"
)

type ordersRepo struct {
	Db *sqlx.DB
}

func NewOrdersRepository(db *sqlx.DB) entities.OrderRepository {
	return &ordersRepo{
		Db: db,
	}
}

func (r *ordersRepo) CreateOrders(req *entities.OrdersReq2) (*entities.OrdersRes2, error) {
	query := `
	WITH create_shipping AS (
		INSERT INTO "orders" ("shipping_address") 
		VALUES ($4)
		RETURNING id, shipping_address
	)
		INSERT INTO "product_order"(
			"order_id",
			"products",
			"qty",
			"price",
			"status",
			"created_at",
			"enable"
		)
		VALUES ((select id from create_shipping), $1, $2, $3, 'placed_order', $5, true)
		RETURNING "id", (select shipping_address from create_shipping), "order_id", "price", "qty", "products", "status", "created_at";
	`

	order := new(entities.OrdersRes2)
	times := time.Now()
	collections := make(map[string][]entities.Product)
	product := []entities.Product{}
	var totalQty, totalPrice int

	for _, v := range req.Product.Item {
		//fmt.Println(v.Id)
		productRes, err := r.QueryCart(v.Id)
		p := entities.Product{Id: productRes.Id, Gender: strings.ToLower(productRes.Gender), StyleType: strings.ToLower(productRes.StyleType), StyleDetail: productRes.StyleDetail, Size: strings.ToLower(productRes.Size), Price: productRes.Price, Qty: v.Qty, TotalPrice: productRes.Price * float64(v.Qty)}
		if err != nil {
			fmt.Println(err.Error())
		}
		totalQty += v.Qty
		totalPrice += int(productRes.Price) * v.Qty
		product = append(product, p)
	}

	collections["item"] = append(collections["item"], product...)

	rows, err := r.Db.Queryx(query, collections, totalQty, totalPrice, req.Shipping, times)

		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		for rows.Next() {
			if err := rows.StructScan(order); err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
		}
	
	return order, nil
}

func (r *ordersRepo) QueryCart(id int)(*entities.Product, error){
	query := `select id, gender, style_type, style_detail, size, price from products where id=$1;`

	product := new(entities.Product)
	//a := entities.Product{}
	if err := r.Db.Get(product, query, id) ; err != nil{
		return nil, err
	}

	return product, nil

}

func (r *ordersRepo) GetOrder(params *entities.QueryParams) (list []*entities.GetOrderRes, err error){
	lists := make([]*entities.GetOrderRes, 0)
	query := `SELECT id, order_id, products, qty, price FROM product_order WHERE enable=true`

	 if params.Sdate != "" && params.Edate != "" {
	 	query += fmt.Sprintf(" AND DATE(created_at) BETWEEN '%v' AND '%v'", params.Sdate, params.Edate)
	 }

	if params.Status != "" {
		query += fmt.Sprintf(" AND status='%v'", strings.ToLower(params.Status))
	}

	pages := params.PerPage * (params.Page - 1)
	query += fmt.Sprintf(` ORDER BY id LIMIT %d OFFSET %d`, params.PerPage, pages)

	rows, err := r.Db.Query(query)

	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		orders := new(entities.GetOrderRes)
		fmt.Println(params.Page)
		err := rows.Scan(&orders.ID, &orders.OrderID, &orders.Products, &orders.Qty, &orders.Price)
		if err != nil{
			return nil, err
		}
		lists = append(lists, orders)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return lists, nil
}