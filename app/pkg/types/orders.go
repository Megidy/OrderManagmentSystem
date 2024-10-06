package types

type Dish struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	OrderId    int    `json:"order_id"`
	CustomerId int    `json:"customer_id"`
	Dishes     []Dish `json:"dishes"`
	Created    string `json:"created"`
	Status     string `json:"status"`
}
type CreateOrder struct {
	Dishes []Dish `json:"dishes"`
}

type CheckOrder struct {
	OrderId int
	Status  string
}
