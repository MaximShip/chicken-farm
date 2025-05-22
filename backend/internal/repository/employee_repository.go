package repository

import (
	"chicken-farm/internal/model"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(employee *model.Employee) error {
	tx := r.db.Begin()

	if err := tx.Create(employee).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, cageID := range employee.Cages {
		employeeCage := model.EmployeeCage{
			EmployeeID: employee.ID,
			CageID:     cageID,
		}
		if err := tx.Create(&employeeCage).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *EmployeeRepository) GetByID(id uint) (*model.Employee, error) {
	var employee model.Employee
	err := r.db.First(&employee, id).Error
	if err != nil {
		return nil, err
	}

	var employeeCages []model.EmployeeCage
	if err := r.db.Where("employee_id = ?", employee.ID).Find(&employeeCages).Error; err != nil {
		return nil, err
	}

	employee.Cages = make([]uint, len(employeeCages))
	for i, ec := range employeeCages {
		employee.Cages[i] = ec.CageID
	}

	return &employee, nil
}

func (r *EmployeeRepository) GetAll() ([]model.Employee, error) {
	var employees []model.Employee
	err := r.db.Find(&employees).Error
	if err != nil {
		return nil, err
	}

	for i := range employees {
		var employeeCages []model.EmployeeCage
		if err := r.db.Where("employee_id = ?", employees[i].ID).Find(&employeeCages).Error; err != nil {
			return nil, err
		}

		employees[i].Cages = make([]uint, len(employeeCages))
		for j, ec := range employeeCages {
			employees[i].Cages[j] = ec.CageID
		}
	}

	return employees, nil
}

func (r *EmployeeRepository) Update(employee *model.Employee) error {
	tx := r.db.Begin()

	if err := tx.Save(employee).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("employee_id = ?", employee.ID).Delete(&model.EmployeeCage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, cageID := range employee.Cages {
		employeeCage := model.EmployeeCage{
			EmployeeID: employee.ID,
			CageID:     cageID,
		}
		if err := tx.Create(&employeeCage).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *EmployeeRepository) Delete(id uint) error {
	tx := r.db.Begin()

	if err := tx.Where("employee_id = ?", id).Delete(&model.EmployeeCage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&model.Employee{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *EmployeeRepository) GetEmployeeChickenCount(employeeID uint) (int, error) {
	var count int64

	err := r.db.Table("employee_cages").
		Joins("JOIN chickens ON employee_cages.cage_id = chickens.cage_id").
		Where("employee_cages.employee_id = ?", employeeID).
		Count(&count).Error

	return int(count), err
}

func (r *EmployeeRepository) GetAllEmployeeChickenCounts() (map[uint]int, error) {
	type Result struct {
		EmployeeID   uint
		ChickenCount int64
	}

	var results []Result
	err := r.db.Table("employee_cages").
		Select("employee_cages.employee_id, COUNT(chickens.id) as chicken_count").
		Joins("JOIN chickens ON employee_cages.cage_id = chickens.cage_id").
		Group("employee_cages.employee_id").
		Scan(&results).Error

	counts := make(map[uint]int)
	for _, r := range results {
		counts[r.EmployeeID] = int(r.ChickenCount)
	}

	return counts, err
}

func (r *EmployeeRepository) GetEmployeeEggCount(employeeID uint, startDate, endDate string) (int, error) {
	var count int64

	err := r.db.Table("employee_cages").
		Joins("JOIN farm_records ON employee_cages.cage_id = farm_records.cage_id").
		Where("employee_cages.employee_id = ? AND farm_records.date BETWEEN ? AND ? AND farm_records.has_egg = true",
			employeeID, startDate, endDate).
		Count(&count).Error

	return int(count), err
}

func (r *EmployeeRepository) GetAllEmployeeEggCounts(startDate, endDate string) (map[uint]int, error) {
	type Result struct {
		EmployeeID uint
		EggCount   int64
	}

	var results []Result
	err := r.db.Table("employee_cages").
		Select("employee_cages.employee_id, COUNT(farm_records.id) as egg_count").
		Joins("JOIN farm_records ON employee_cages.cage_id = farm_records.cage_id").
		Where("farm_records.date BETWEEN ? AND ? AND farm_records.has_egg = true", startDate, endDate).
		Group("employee_cages.employee_id").
		Scan(&results).Error

	counts := make(map[uint]int)
	for _, r := range results {
		counts[r.EmployeeID] = int(r.EggCount)
	}

	return counts, err
}
