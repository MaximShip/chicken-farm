package service

import (
	"chicken-farm/internal/model"
	"chicken-farm/internal/repository"
)

type ReportService struct {
	chickenRepo  *repository.ChickenRepository
	employeeRepo *repository.EmployeeRepository
	farmRepo     *repository.FarmRepository
}

func NewReportService(
	chickenRepo *repository.ChickenRepository,
	employeeRepo *repository.EmployeeRepository,
	farmRepo *repository.FarmRepository,
) *ReportService {
	return &ReportService{
		chickenRepo:  chickenRepo,
		employeeRepo: employeeRepo,
		farmRepo:     farmRepo,
	}
}

func (s *ReportService) GetAvgEggsByWeightAndAge(weight float64, age int) (float64, error) {
	return s.chickenRepo.GetAvgEggsByWeightAndAge(weight, age)
}

type EggStats struct {
	TotalEggs int     `json:"total_eggs"`
	TotalCost float64 `json:"total_cost"`
}

func (s *ReportService) GetTotalEggStats(startDate, endDate string) (*EggStats, error) {
	count, err := s.farmRepo.GetEggCountByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	cost, err := s.farmRepo.GetTotalEggCost(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &EggStats{
		TotalEggs: count,
		TotalCost: cost,
	}, nil
}

type EmployeeEggStats struct {
	EmployeeID   uint   `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	EggCount     int    `json:"egg_count"`
}

func (s *ReportService) GetEmployeeEggStats(startDate, endDate string) ([]EmployeeEggStats, error) {
	eggCounts, err := s.employeeRepo.GetAllEmployeeEggCounts(startDate, endDate)
	if err != nil {
		return nil, err
	}

	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	stats := make([]EmployeeEggStats, 0, len(employees))
	for _, employee := range employees {
		eggCount, ok := eggCounts[employee.ID]
		if !ok {
			eggCount = 0
		}

		stats = append(stats, EmployeeEggStats{
			EmployeeID:   employee.ID,
			EmployeeName: employee.FullName,
			EggCount:     eggCount,
		})
	}

	return stats, nil
}

func (s *ReportService) GetLowProductivityChickens() ([]model.Chicken, error) {
	return s.chickenRepo.GetChickensWithLowProductivity()
}

type MostProductiveChickenStats struct {
	ChickenID   uint `json:"chicken_id"`
	CageID      uint `json:"cage_id"`
	CageNumber  int  `json:"cage_number"`
	EggPerMonth int  `json:"egg_per_month"`
}

func (s *ReportService) GetMostProductiveChickenStats() (*MostProductiveChickenStats, error) {
	chicken, err := s.chickenRepo.GetMostProductiveChicken()
	if err != nil {
		return nil, err
	}

	cage, err := s.farmRepo.GetCageByID(chicken.CageID)
	if err != nil {
		return nil, err
	}

	return &MostProductiveChickenStats{
		ChickenID:   chicken.ID,
		CageID:      chicken.CageID,
		CageNumber:  cage.Number,
		EggPerMonth: chicken.EggPerMonth,
	}, nil
}

type EmployeeChickenCountStats struct {
	EmployeeID   uint   `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	ChickenCount int    `json:"chicken_count"`
}

func (s *ReportService) GetEmployeeChickenCountStats() ([]EmployeeChickenCountStats, error) {
	chickenCounts, err := s.employeeRepo.GetAllEmployeeChickenCounts()
	if err != nil {
		return nil, err
	}

	employees, err := s.employeeRepo.GetAll()
	if err != nil {
		return nil, err
	}

	stats := make([]EmployeeChickenCountStats, 0, len(employees))
	for _, employee := range employees {
		chickenCount, ok := chickenCounts[employee.ID]
		if !ok {
			chickenCount = 0
		}

		stats = append(stats, EmployeeChickenCountStats{
			EmployeeID:   employee.ID,
			EmployeeName: employee.FullName,
			ChickenCount: chickenCount,
		})
	}

	return stats, nil
}
