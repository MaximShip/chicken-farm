package service

import (
	"errors"

	"chicken-farm/internal/model"
	"chicken-farm/internal/repository"
)

type ChickenService struct {
	chickenRepo *repository.ChickenRepository
	farmRepo    *repository.FarmRepository
}

func NewChickenService(
	chickenRepo *repository.ChickenRepository,
	farmRepo *repository.FarmRepository,
) *ChickenService {
	return &ChickenService{
		chickenRepo: chickenRepo,
		farmRepo:    farmRepo,
	}
}

func (s *ChickenService) CreateChicken(chicken *model.Chicken) error {
	_, err := s.farmRepo.GetCageByID(chicken.CageID)
	if err != nil {
		return errors.New("cage not found")
	}

	existingChicken, err := s.chickenRepo.GetByCageID(chicken.CageID)
	if err == nil && existingChicken != nil {
		return errors.New("cage is already occupied by another chicken")
	}

	return s.chickenRepo.Create(chicken)
}

func (s *ChickenService) GetChickenByID(id uint) (*model.Chicken, error) {
	return s.chickenRepo.GetByID(id)
}

func (s *ChickenService) GetAllChickens() ([]model.Chicken, error) {
	return s.chickenRepo.GetAll()
}

func (s *ChickenService) UpdateChicken(chicken *model.Chicken) error {
	_, err := s.chickenRepo.GetByID(chicken.ID)
	if err != nil {
		return errors.New("chicken not found")
	}

	oldChicken, err := s.chickenRepo.GetByID(chicken.ID)
	if err != nil {
		return err
	}

	if oldChicken.CageID != chicken.CageID {
		_, err := s.farmRepo.GetCageByID(chicken.CageID)
		if err != nil {
			return errors.New("new cage not found")
		}

		existingChicken, err := s.chickenRepo.GetByCageID(chicken.CageID)
		if err == nil && existingChicken != nil && existingChicken.ID != chicken.ID {
			return errors.New("new cage is already occupied by another chicken")
		}
	}

	return s.chickenRepo.Update(chicken)
}

func (s *ChickenService) DeleteChicken(id uint) error {
	_, err := s.chickenRepo.GetByID(id)
	if err != nil {
		return errors.New("chicken not found")
	}

	return s.chickenRepo.Delete(id)
}

func (s *ChickenService) GetAvgEggsByWeightAndAge(weight float64, age int) (float64, error) {
	return s.chickenRepo.GetAvgEggsByWeightAndAge(weight, age)
}

func (s *ChickenService) GetChickensWithLowProductivity() ([]model.Chicken, error) {
	return s.chickenRepo.GetChickensWithLowProductivity()
}

func (s *ChickenService) GetMostProductiveChicken() (*model.Chicken, error) {
	return s.chickenRepo.GetMostProductiveChicken()
}

func (s *ChickenService) GetCageWithMostEggs() (*model.Cage, error) {
	cageID, err := s.farmRepo.GetCageWithMostEggs()
	if err != nil {
		return nil, err
	}

	return s.farmRepo.GetCageByID(cageID)
}
