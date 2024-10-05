package db

import (
	"database/sql"

	"github.com/Megidy/OrderManagmentSystem/pkg/config"
	"github.com/Megidy/OrderManagmentSystem/pkg/types"
)

var db *sql.DB

func init() {
	config.ConnectDB()
	db = config.GetDB()
}
func CreateOrder(order types.Order) error {

	_, err := db.Exec("INSERT INTO orders (customer_id, order_id, status) VALUES (?, ?, ?)",
		order.CustomerId, order.OrderId, order.Status)
	if err != nil {
		return err
	}

	for _, dish := range order.Dishes {
		err := CreateDish(order.OrderId, dish)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateDish(orderId int, dish types.Dish) error {

	_, err := db.Exec("INSERT INTO dishes (order_id, name, quantity) VALUES (?, ?, ?)",
		orderId, dish.Name, dish.Quantity)
	if err != nil {
		return err
	}
	return nil
}
func CheckOrdersStatus() ([]types.CheckOrder, error) {
	var Orders []types.CheckOrder
	query, err := db.Query("select status,order_id from orders")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	for query.Next() {
		var checkOrder types.CheckOrder
		err = query.Scan(&checkOrder.Status, &checkOrder.OrderId)
		if err != nil {
			return nil, err
		}
		Orders = append(Orders, checkOrder)
	}

	return Orders, nil

}
func ChangeOrderStatus(order types.Order, status string) error {
	_, err := db.Exec("update orders set status =? where order_id=? and customer_id=?", status, order.OrderId, order.CustomerId)
	if err != nil {
		return err
	}
	return nil
}
