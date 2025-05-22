package service

import (
	"errors"

	"chicken-farm/internal/model"
	"chicken-farm/internal/repository"
)

type EmployeeService struct {
	employeeRepo *repository.EmployeeRepository
	farmRepo     *repository.FarmRepository
}

func NewEmployeeService(
	employeeRepo *repository.EmployeeRepository,
	farmRepo *repository.FarmRepository,
) *EmployeeService {
	return &EmployeeService{
		employeeRepo: employeeRepo,
		farmRepo:     farmRepo,
	}
}

// CreateEmplyee создает нового работника
func (s *EmployeeService) CreateEmployee(employee *model.Employee) error {
	// Проверяем, что все клетки существуют
	for _, cageID := range employee.Cages {
		_, err := s.farmRepo.GetCageByID(cageID)
		if err != nil {
			return errors.New("cage not found")
		}
	}

	return s.employeeRepo.Create(employee)
}

func (s *EmployeeService) GetEmployeeByID(id uint) (*model.Employee, error) {
	return s.employeeRepo.GetByID(id)
}

func (s *EmployeeService) GetAllEmployees() ([]model.Employee, error) {
	return s.employeeRepo.GetAll()
}

func (s *EmployeeService) UpdateEmployee(employee *model.Employee) error {
	_, err := s.employeeRepo.GetByID(employee.ID)
	if err != nil {
		return errors.New("employee not found")
	}

	for _, cageID := range employee.Cages {
		_, err := s.farmRepo.GetCageByID(cageID)
		if err != nil {
			return errors.New("cage not found")
		}
	}

	return s.employeeRepo.Update(employee)
}

func (s *EmployeeService) DeleteEmployee(id uint) error {
	_, err := s.employeeRepo.GetByID(id)
	if err != nil {
		return errors.New("employee not found")
	}

	return s.employeeRepo.Delete(id)
}

func (s *EmployeeService) GetEmployeeChickenCount(employeeID uint) (int, error) {
	_, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return 0, errors.New("employee not found")
	}

	return s.employeeRepo.GetEmployeeChickenCount(employeeID)
}

func (s *EmployeeService) GetAllEmployeeChickenCounts() (map[uint]int, error) {
	return s.employeeRepo.GetAllEmployeeChickenCounts()
}

func (s *EmployeeService) GetEmployeeEggCount(employeeID uint, startDate, endDate string) (int, error) {
	_, err := s.employeeRepo.GetByID(employeeID)
	if err != nil {
		return 0, errors.New("employee not found")
	}

	return s.employeeRepo.GetEmployeeEggCount(employeeID, startDate, endDate)
}

func (s *EmployeeService) GetAllEmployeeEggCounts(startDate, endDate string) (map[uint]int, error) {
	return s.employeeRepo.GetAllEmployeeEggCounts(startDate, endDate)
}
