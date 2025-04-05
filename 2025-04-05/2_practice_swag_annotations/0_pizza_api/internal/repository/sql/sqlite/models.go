package sqlite

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	FirstName string         `gorm:"type:text;not null"`
	LastName  string         `gorm:"type:text;not null"`
	Address   string         `gorm:"type:text;not null"`
	Email     string         `gorm:"type:text;not null;unique"`
	Orders    []Order        `gorm:"foreignKey:CustomerID"`
}

type Order struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Type       string         `gorm:"type:text;not null"`
	Size       string         `gorm:"type:text;not null"`
	Quantity   int            `gorm:"type:integer;not null"`
	CustomerID uint           `gorm:"type:integer;not null;index"`
	Status     string         `gorm:"type:text;not null"`
	Customer   Customer       `gorm:"foreignKey:CustomerID"`
}
