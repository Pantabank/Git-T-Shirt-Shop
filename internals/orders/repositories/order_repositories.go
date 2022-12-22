package repositories

import (
	//"encoding/json"

	"fmt"
	//"time"

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
		INSERT INTO "product_order"(
			"order_id",
			"products",
			"qty",
			"price"
		)
		VALUES ($1, $2, $3, $4)
		RETURNING "id", "order_id", "price", "qty", "products";
	`

	order := new(entities.OrdersRes2)
	//times := time.Now()

	collections := make(map[string][]entities.Product)
	product := []entities.Product{}
	var totalQty, totalPrice int

	for _, v := range req.Product.Item {
		//fmt.Println(v.Id)
		productRes, err := r.QueryCart(v.Id)
		p := entities.Product{Id: productRes.Id, Gender: productRes.Gender, StyleType: productRes.StyleType, StyleDetail: productRes.StyleDetail, Size: productRes.Size, Price: productRes.Price, Qty: v.Qty, TotalPrice: productRes.Price * float64(v.Qty)}
		if err != nil {
			fmt.Println(err.Error())
		}
		totalQty += v.Qty
		totalPrice += int(productRes.Price) * v.Qty
		product = append(product, p)
	}

	collections["item"] = append(collections["item"], product...)

	rows, err := r.Db.Queryx(query, 142, collections, totalQty, totalPrice)

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

func (r *ordersRepo) GetOrder(params *entities.QueryParams) (list []*entities.OrdersRes, err error){
	lists := make([]*entities.OrdersRes, 0)
	query := `SELECT id, product_id, gender, style_type, style_detail, size, price, shipping_address, status, created_at FROM orders WHERE enable=true`

	if params.Sdate != "" && params.Edate != "" {
		query += fmt.Sprintf(" AND DATE(created_at) BETWEEN '%v' AND '%v'", params.Sdate, params.Edate)
	}

	if params.Status != "" {
		query += fmt.Sprintf(" AND status='%v'", params.Status)
	}

	pages := params.PerPage * (params.Page - 1)
	query += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.PerPage, pages)

	rows, err := r.Db.Query(query)

	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		orders := new(entities.OrdersRes)
		fmt.Println(params.Page)
		err := rows.Scan(&orders.Id, &orders.ProductID, &orders.Gender, &orders.StyleType, &orders.StyleDetail, &orders.Size, &orders.Price, &orders.ShippingAddress, &orders.Status, &orders.CreatedAt)
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