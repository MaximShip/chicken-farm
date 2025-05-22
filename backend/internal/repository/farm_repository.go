package repository

import (
	"chicken-farm/internal/model"

	"gorm.io/gorm"
)

type FarmRepository struct {
	db *gorm.DB
}

// NewFarmRepository создает новый репозиторий для птицефабрики
func NewFarmRepository(db *gorm.DB) *FarmRepository {
	return &FarmRepository{db: db}
}

func (r *FarmRepository) Create(farm *model.Farm) error {
	return r.db.Create(farm).Error
}

func (r *FarmRepository) GetByID(id uint) (*model.Farm, error) {
	var farm model.Farm
	err := r.db.First(&farm, id).Error
	if err != nil {
		return nil, err
	}
	return &farm, nil
}

func (r *FarmRepository) GetByDateRange(startDate, endDate string) ([]model.Farm, error) {
	var farms []model.Farm
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&farms).Error
	return farms, err
}

func (r *FarmRepository) GetByChickenID(chickenID uint) ([]model.Farm, error) {
	var farms []model.Farm
	err := r.db.Where("chicken_id = ?", chickenID).Find(&farms).Error
	return farms, err
}

func (r *FarmRepository) GetEggCountByDateRange(startDate, endDate string) (int, error) {
	var count int64
	err := r.db.Model(&model.Farm{}).
		Where("date BETWEEN ? AND ? AND has_egg = true", startDate, endDate).
		Count(&count).Error

	return int(count), err
}

func (r *FarmRepository) GetTotalEggCost(startDate, endDate string) (float64, error) {
	eggCount, err := r.GetEggCountByDateRange(startDate, endDate)
	if err != nil {
		return 0, err
	}

	return float64(eggCount) * model.GetEggPrice(), nil
}

func (r *FarmRepository) GetCageWithMostEggs() (uint, error) {
	type Result struct {
		CageID   uint
		EggCount int64
	}

	var result Result
	err := r.db.Model(&model.Farm{}).
		Select("cage_id, COUNT(*) as egg_count").
		Where("has_egg = true").
		Group("cage_id").
		Order("egg_count DESC").
		First(&result).Error

	return result.CageID, err
}

func (r *FarmRepository) CreateCage(cage *model.Cage) error {
	return r.db.Create(cage).Error
}

func (r *FarmRepository) GetCageByID(id uint) (*model.Cage, error) {
	var cage model.Cage
	err := r.db.First(&cage, id).Error
	if err != nil {
		return nil, err
	}
	return &cage, nil
}

func (r *FarmRepository) GetAllCages() ([]model.Cage, error) {
	var cages []model.Cage
	err := r.db.Find(&cages).Error
	return cages, err
}

func (r *FarmRepository) GetEmptyCages() ([]model.Cage, error) {
	var allCageIDs []uint
	var occupiedCageIDs []uint

	if err := r.db.Model(&model.Cage{}).Pluck("id", &allCageIDs).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&model.Chicken{}).Pluck("cage_id", &occupiedCageIDs).Error; err != nil {
		return nil, err
	}

	emptyCageIDs := make([]uint, 0)
	for _, cageID := range allCageIDs {
		found := false
		for _, occupiedID := range occupiedCageIDs {
			if cageID == occupiedID {
				found = true
				break
			}
		}
		if !found {
			emptyCageIDs = append(emptyCageIDs, cageID)
		}
	}

	var emptyCages []model.Cage
	if len(emptyCageIDs) > 0 {
		err := r.db.Where("id IN ?", emptyCageIDs).Find(&emptyCages).Error
		return emptyCages, err
	}

	return emptyCages, nil
}

func (r *FarmRepository) UpdateConfigParam(key, value string) error {
	var param model.ConfigParam

	err := r.db.Where("key = ?", key).First(&param).Error
	if err == gorm.ErrRecordNotFound {
		param = model.ConfigParam{
			Key:   key,
			Value: value,
		}
		return r.db.Create(&param).Error
	} else if err != nil {
		return err
	}

	param.Value = value
	return r.db.Save(&param).Error
}

func (r *FarmRepository) GetConfigParam(key string) (string, error) {
	var param model.ConfigParam
	err := r.db.Where("key = ?", key).First(&param).Error
	if err != nil {
		return "", err
	}
	return param.Value, nil
}
