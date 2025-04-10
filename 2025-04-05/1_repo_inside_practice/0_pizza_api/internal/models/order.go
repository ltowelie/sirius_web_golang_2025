package models

import "time"

type PizzaSize string

// Пусть пока будет захардкожено для примера
const (
	Small  PizzaSize = "small"
	Medium PizzaSize = "medium"
	Large  PizzaSize = "large"
)

type PizzaType string

// Пусть пока будет захардкожено для примера
const (
	Margarita PizzaType = "margarita"
	Marinara  PizzaType = "marinara"
	Pepperoni PizzaType = "pepperoni"
)

type OrderStatus string

const (
	New        OrderStatus = "new"
	Preparing  OrderStatus = "preparing"
	Delivering OrderStatus = "delivering"
	Delivered  OrderStatus = "delivered"
	Cancelled  OrderStatus = "cancelled"
)

// Order А с этой моделью как раз и будем работать
type Order struct {
	ID         int         `json:"id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	DeletedAt  time.Time   `json:"deleted_at"`
	Type       PizzaType   `json:"type"`
	Size       PizzaSize   `json:"size"`
	Quantity   int         `json:"quantity"`
	CustomerID int         `json:"customer_id"`
	Status     OrderStatus `json:"status"`
}

type Customer struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
}
