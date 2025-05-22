package model

import (
	"time"
)

type Chicken struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CageID      uint      `json:"cage_id" gorm:"not null"`
	Weight      float64   `json:"weight" gorm:"not null"`        // вес в килограммах
	Age         int       `json:"age" gorm:"not null"`           // возраст в месяцах
	EggPerMonth int       `json:"egg_per_month" gorm:"not null"` // количество яиц в месяц
	Breed       string    `json:"breed" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Chicken) TableName() string {
	return "chickens"
}
