package types

import "time"

type OrderServices interface {
	CreateOrder(Order) error
	GetOrder(Order) (Order, error)
	SendOrder(Order) error
}

type Dish struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	CustomersId int       `json:"id"`
	Dishes      []Dish    `json:"dishes"`
	Created     time.Time `json:"created"`
	Status      string    `json:"status"`
}
