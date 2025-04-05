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
	ID         int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	Type       PizzaType
	Size       PizzaSize
	Quantity   int
	CustomerID int
	Status     OrderStatus
}

type Customer struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	FirstName string
	LastName  string
	Address   string
	Email     string
}
