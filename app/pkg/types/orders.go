package types

import (
	"time"
)

type Dish struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	CustomerId int       `json:"customer_id"`
	OrderId    int       `json:"order_id"`
	Dishes     []Dish    `json:"dishes"`
	Created    time.Time `json:"created"`
	Status     string    `json:"status"`
}
type CreateOrder struct {
	Dishes []Dish `json:"dishes"`
	Status string `json:"status"`
}
