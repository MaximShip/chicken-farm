package model

import (
	"time"
)

type Employee struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FullName     string    `json:"full_name" gorm:"not null"`
	PassportData string    `json:"passport_data" gorm:"not null;uniqueIndex"`
	Salary       float64   `json:"salary" gorm:"not null"`
	Cages        []uint    `json:"cages" gorm:"-"` // Список ID клеток
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EmployeeCage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	EmployeeID uint      `json:"employee_id" gorm:"not null"`
	CageID     uint      `json:"cage_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Employee) TableName() string {
	return "employees"
}

func (EmployeeCage) TableName() string {
	return "employee_cages"
}
