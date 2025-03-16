package models

// Пакет также может называться Domain, Entity

type Pizza struct {
	ID          int
	Name        string
	Description string
	// Тут конечно нужно уточнение - хранить лучше в каких то единицах измерения или перечислениях типа "L", "M", "S"
	Size string
}
