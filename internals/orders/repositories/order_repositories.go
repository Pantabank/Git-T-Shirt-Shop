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
		INSERT INTO "product_order"("order_id", product_id, "products")
		VALUES ($1, $2, $3)
		RETURNING "id", (select shipping_address from orders where id=$1 ), "order_id";
	`

	orderid, err := r.GetOrderID(req.Shipping)
	if err != nil {
		fmt.Println(err.Error())
	}
	times := time.Now()
	order := new(entities.OrdersRes2)
	product := []entities.Product{}
	collections := make(map[string][]entities.Product)
	var totalQty, totalPrice int

	for _, v := range req.Product.Item {
		productRes, err := r.QueryCart(v.Id)
		p := entities.Product{Id: productRes.Id, Gender: strings.ToLower(productRes.Gender), StyleType: strings.ToLower(productRes.StyleType), StyleDetail: productRes.StyleDetail, Size: strings.ToLower(productRes.Size), Price: productRes.Price, Qty: v.Qty, TotalPrice: productRes.Price * float64(v.Qty)}
		if err != nil {
			fmt.Println(err.Error())
		}
		totalQty += v.Qty
		totalPrice += int(productRes.Price) * v.Qty
		product = append(product, p)

		rows, err := r.Db.Queryx(query, orderid.Id, productRes.Id, p)

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
	}

	collections["item"] = append(collections["item"], product...)
	r.UpdateOrders(totalQty, &totalPrice, order.OrderID) 
	pd := entities.OrdersRes2{Id: order.Id, OrderID: order.OrderID, Qty: totalQty, Price: totalPrice, Shipping:  order.Shipping, Product: entities.ProductItem{product}, Status: "placed_order", CreatedAt:times  }
	fmt.Println(collections)
	return &pd, nil
}

func (r *ordersRepo) QueryCart(id int)(*entities.Product, error){
	query := `select id, gender, style_type, style_detail, size, price from products where id=$1;`
	product := new(entities.Product)
	if err := r.Db.Get(product, query, id) ; err != nil{
		return nil, err
	}
	return product, nil
}

func (r *ordersRepo) GetOrderID(address *entities.Shipping)(*entities.AddressesRes, error){
	query := `
		INSERT INTO "orders" ("shipping_address", "enable") 
		VALUES ($1, true)
		RETURNING id, shipping_address
	`

	addresses := new(entities.AddressesRes)
	rows, err := r.Db.Queryx(query, address)

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

func (r *ordersRepo) UpdateOrders(qty int, price *int, order_id int) error {
	query := `
		UPDATE orders SET qty=$1, price=$2, status=$3, created_at=$4 WHERE id=$5
	`
	times := time.Now()
	fmt.Println(order_id)
	_, err := r.Db.Exec(query, qty, price, "placed_order", times, order_id)
	if err != nil{
		panic(err)
	}
	return err

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