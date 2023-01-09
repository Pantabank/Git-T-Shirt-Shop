package repositories

import (
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

func (r *ordersRepo) AddOrders(order_id, product_id int, product *entities.Product) error {
	query := `
		INSERT INTO "product_order"("order_id", product_id, "products")
		VALUES ($1, $2, $3)
		RETURNING "id", (select shipping_address from orders where id=$1 ), "order_id";
	`

	rows, err := r.Db.Queryx(query, order_id, product_id, product)

		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		for rows.Next() {
			if err := rows.StructScan(new(entities.OrdersRes2)); err != nil {
				fmt.Println(err.Error())
				return nil
			}
		}

		return nil
}

func (r *ordersRepo) QueryCart(id int)(*entities.Product, error){
	query := `select id, gender, style_type, style_detail, size, price from products where id=$1;`
	product := new(entities.Product)
	if err := r.Db.Get(product, query, id) ; err != nil{
		return nil, err
	}
	return product, nil
}

func (r *ordersRepo) GetOrderID(address *entities.Shipping, uid any)(*entities.AddressesRes, error){

	query := `
		INSERT INTO "orders" ("shipping_address", "enable", customer_id) 
		VALUES ($1, true, $2)
		RETURNING id, shipping_address
	`

	addresses := new(entities.AddressesRes)
	rows, err := r.Db.Queryx(query, address, uid)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	for rows.Next() {
		if err := rows.StructScan(addresses); err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return addresses, nil
}

func (r *ordersRepo) UpdateOrders(qty int, price int, order_id int) error {
	query := `
		UPDATE orders SET qty=$1, price=$2, status=$3, created_at=$4 WHERE id=$5
	`
	times := time.Now()
	_, err := r.Db.Exec(query, qty, price, "placed_order", times, order_id)
	if err != nil{
		panic(err)
	}
	return err

}

func (r *ordersRepo) GetOrder(params *entities.QueryParams) (list []*entities.GetOrderRes, err error){
	lists := make([]*entities.GetOrderRes, 0)
	query := `SELECT o.id, o.shipping_address, o.qty, o.price, o.status, o.created_at AS shipping_address, concat('{"item":', jsonb_agg(po.products), '}')::jsonb AS orders 
	FROM orders o 
	LEFT JOIN product_order po ON o.id = po.order_id
	WHERE enable=true
	`

	 if params.Sdate != "" && params.Edate != "" {
	 	query += fmt.Sprintf(" AND DATE(o.created_at) BETWEEN '%v' AND '%v'", params.Sdate, params.Edate)
	 }

	if params.Status != "" {
		query += fmt.Sprintf(" AND o.status='%v'", strings.ToLower(params.Status))
	}

	pages := params.PerPage * (params.Page - 1)
	query += fmt.Sprintf(` GROUP BY o.id ORDER BY id LIMIT %d OFFSET %d`, params.PerPage, pages)

	rows, err := r.Db.Query(query)

	if err != nil{
		return nil, err
	}
	defer rows.Close()


	for rows.Next(){
		orders := new(entities.GetOrderRes)
		err := rows.Scan(&orders.ID, &orders.Shipping, &orders.Qty, &orders.Price, &orders.Status, &orders.CreatedAt, &orders.Product)
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