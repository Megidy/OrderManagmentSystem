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
	// Спочатку вставляємо замовлення в таблицю orders
	_, err := db.Exec("INSERT INTO orders (customer_id, order_id, status) VALUES (?, ?, ?)",
		order.CustomerId, order.OrderId, order.Status)
	if err != nil {
		return err
	}

	// Тепер вставляємо кожну страву з масиву Dishes в таблицю dishes
	for _, dish := range order.Dishes {
		err := CreateDish(order.OrderId, dish)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateDish(orderId int, dish types.Dish) error {
	// Вставляємо кожну страву в таблицю dishes з відповідним order_id
	_, err := db.Exec("INSERT INTO dishes (order_id, name, quantity) VALUES (?, ?, ?)",
		orderId, dish.Name, dish.Quantity)
	if err != nil {
		return err
	}
	return nil
}
