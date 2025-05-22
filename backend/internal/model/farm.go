package model

import (
	"time"
)

type Farm struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Date      time.Time `json:"date" gorm:"not null"`
	CageID    uint      `json:"cage_id" gorm:"not null"`
	ChickenID uint      `json:"chicken_id" gorm:"not null"`
	HasEgg    bool      `json:"has_egg" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Cage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Number    int       `json:"number" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Farm) TableName() string {
	return "farm_records"
}

func (Cage) TableName() string {
	return "cages"
}

type ConfigParam struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"not null;uniqueIndex"`
	Value     string    `json:"value" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ConfigParam) TableName() string {
	return "config_params"
}

func GetEggPrice() float64 {
	return 10.0
}
