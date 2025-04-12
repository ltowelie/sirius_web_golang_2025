package models

import "time"

// PizzaSize определяет размер пиццы
// @Description Размер заказываемой пиццы
type PizzaSize string

// Размеры пиццы
// @Enum small medium large
const (
	Small  PizzaSize = "small"  // Маленькая (25 см)
	Medium PizzaSize = "medium" // Средняя (30 см)
	Large  PizzaSize = "large"  // Большая (35 см)
)

// PizzaType определяет тип пиццы
// @Description Доступные виды пиццы
type PizzaType string

// Типы пиццы
// @Enum margarita marinara pepperoni
const (
	Margarita PizzaType = "margarita" // Маргарита
	Marinara  PizzaType = "marinara"  // Маринара
	Pepperoni PizzaType = "pepperoni" // Пепперони
)

// OrderStatus определяет статус заказа
// @Description Статус выполнения заказа
type OrderStatus string

// Статусы заказа
// @Enum new preparing delivering delivered cancelled
const (
	New        OrderStatus = "new"        // Новый заказ
	Preparing  OrderStatus = "preparing"  // Готовится
	Delivering OrderStatus = "delivering" // Доставляется
	Delivered  OrderStatus = "delivered"  // Доставлен
	Cancelled  OrderStatus = "cancelled"  // Отменен
)

// Order представляет объект заказа
// @Description Модель заказа пиццы
type Order struct {
	ID         int         `json:"id" example:"1"`                            // Уникальный идентификатор заказа
	CreatedAt  time.Time   `json:"created_at" example:"2025-01-01T12:00:00Z"` // Время создания заказа
	UpdatedAt  time.Time   `json:"updated_at" example:"2025-01-01T12:30:00Z"` // Время последнего обновления
	DeletedAt  time.Time   `json:"deleted_at,omitempty"`                      // Время удаления
	Type       PizzaType   `json:"type" example:"pepperoni"`                  // Тип пиццы
	Size       PizzaSize   `json:"size" example:"medium"`                     // Размер пиццы
	Quantity   int         `json:"quantity" example:"2"`                      // Количество пицц в заказе
	CustomerID int         `json:"customer_id" example:"42"`                  // ID клиента
	Status     OrderStatus `json:"status" example:"new"`                      // Статус заказа
}
