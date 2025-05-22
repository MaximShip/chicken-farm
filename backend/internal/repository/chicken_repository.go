package repository

import (
	"chicken-farm/internal/model"

	"gorm.io/gorm"
)

type ChickenRepository struct {
	db *gorm.DB
}

func NewChickenRepository(db *gorm.DB) *ChickenRepository {
	return &ChickenRepository{db: db}
}

func (r *ChickenRepository) Create(chicken *model.Chicken) error {
	return r.db.Create(chicken).Error
}

func (r *ChickenRepository) GetByID(id uint) (*model.Chicken, error) {
	var chicken model.Chicken
	err := r.db.First(&chicken, id).Error
	if err != nil {
		return nil, err
	}
	return &chicken, nil
}

func (r *ChickenRepository) GetByCageID(cageID uint) (*model.Chicken, error) {
	var chicken model.Chicken
	err := r.db.Where("cage_id = ?", cageID).First(&chicken).Error
	if err != nil {
		return nil, err
	}
	return &chicken, nil
}

func (r *ChickenRepository) GetAll() ([]model.Chicken, error) {
	var chickens []model.Chicken
	err := r.db.Find(&chickens).Error
	return chickens, err
}

func (r *ChickenRepository) GetByWeightAndAge(weight float64, age int) ([]model.Chicken, error) {
	var chickens []model.Chicken
	err := r.db.Where("weight = ? AND age = ?", weight, age).Find(&chickens).Error
	return chickens, err
}

func (r *ChickenRepository) Update(chicken *model.Chicken) error {
	return r.db.Save(chicken).Error
}

func (r *ChickenRepository) Delete(id uint) error {
	return r.db.Delete(&model.Chicken{}, id).Error
}

func (r *ChickenRepository) GetAvgEggsByWeightAndAge(weight float64, age int) (float64, error) {
	var avgEggs float64
	err := r.db.Model(&model.Chicken{}).
		Select("AVG(egg_per_month) as avg_eggs").
		Where("weight = ? AND age = ?", weight, age).
		Scan(&avgEggs).Error

	return avgEggs, err
}

func (r *ChickenRepository) GetChickensWithLowProductivity() ([]model.Chicken, error) {
	var avgEggPerMonth float64
	err := r.db.Model(&model.Chicken{}).Select("AVG(egg_per_month)").Scan(&avgEggPerMonth).Error
	if err != nil {
		return nil, err
	}

	var chickens []model.Chicken
	err = r.db.Where("egg_per_month < ?", avgEggPerMonth).Find(&chickens).Error
	return chickens, err
}

func (r *ChickenRepository) GetMostProductiveChicken() (*model.Chicken, error) {
	var chicken model.Chicken
	err := r.db.Order("egg_per_month DESC").First(&chicken).Error
	if err != nil {
		return nil, err
	}
	return &chicken, nil
}
