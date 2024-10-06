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

	_, err := db.Exec("insert into orders (customer_id, order_id, status) values (?, ?, ?)",
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

func GetOrder(orderId, customerId int) (types.Order, error) {
	var order types.Order
	row := db.QueryRow("select * from orders where order_id=?", orderId)
	err := row.Scan(&order.OrderId, &order.CustomerId, &order.Created, &order.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Order{}, nil
		}
		return types.Order{}, err
	}
	return order, nil

}
func DeleteDishes(orderId int) error {
	_, err := db.Exec("delete from dishes where order_id =?", orderId)

	if err != nil {
		return err
	}
	return nil
}
func DeleteOrder(orderId int) error {
	_, err := db.Exec("delete from orders where order_id =?", orderId)
	if err != nil {
		return err
	}
	err = DeleteDishes(orderId)
	if err != nil {
		return err
	}
	return nil
}

func GetDishes(orderId int) ([]types.Dish, error) {
	var dishes []types.Dish
	query, err := db.Query("select name,quantity from dishes where order_id=?", orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err

	}
	for query.Next() {
		var dish types.Dish
		err = query.Scan(&dish.Name, &dish.Quantity)
		if err != nil {
			return nil, err
		}
		dishes = append(dishes, dish)
	}
	return dishes, nil
}

func CreateDish(orderId int, dish types.Dish) error {

	_, err := db.Exec("insert into dishes (order_id, name, quantity) values (?, ?, ?)",
		orderId, dish.Name, dish.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func CheckcustomersOrdersStatus(orderId int) ([]types.CheckOrder, error) {
	var Orders []types.CheckOrder
	query, err := db.Query("select status,order_id from orders where order_id=?", orderId)
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
