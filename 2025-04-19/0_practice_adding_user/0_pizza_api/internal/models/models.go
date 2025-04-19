package models

type FilterOptions struct {
	Status     OrderStatus
	CustomerID int
}

type Pagination struct {
	Limit  int
	Offset int
}
