package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

// Oil methods

func (s *Storage) AddOil(oil models.OilInput) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	newOil := models.Oil{
		Id:       utils.Id(),
		OilInput: oil,
	}

	storage.Oil = append(storage.Oil, newOil)
	return s.saveStorage(storage)
}

func (s *Storage) GetOil() ([]models.Oil, error) {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return nil, err
	}

	return storage.Oil, nil
}

func (s *Storage) GetOilEntry(id int64) (models.Oil, error) {
	oils, err := s.GetOil()
	if err != nil {
		return models.Oil{}, err
	}

	for _, oil := range oils {
		if oil.Id == id {
			return oil, nil
		}
	}

	return models.Oil{}, nil
}

func (s *Storage) UpdateOil(oil models.Oil) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	for i := range storage.Oil {
		if storage.Oil[i].Id == oil.Id {
			storage.Oil[i] = oil
			break
		}
	}

	return s.saveStorage(storage)
}

func (s *Storage) DeleteOil(id int64) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	for i := range storage.Oil {
		if storage.Oil[i].Id == id {
			storage.Oil = slices.Delete(storage.Oil, i, i+1)
			break
		}
	}

	return s.saveStorage(storage)
}

// OilFillLevel methods

func (s *Storage) AddOilFillLevel(oilFillLevel models.OilFillLevelInput) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	newOilFillLevel := models.OilFillLevel{
		Id:                utils.Id(),
		OilFillLevelInput: oilFillLevel,
	}

	storage.OilFillLevels = append(storage.OilFillLevels, newOilFillLevel)
	return s.saveStorage(storage)
}

func (s *Storage) GetOilFillLevels() ([]models.OilFillLevel, error) {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return nil, err
	}

	for i := range storage.OilFillLevels {
		storage.OilFillLevels[i].Percentage = calculateOilFillPercentage(
			storage.OilFillLevels[i].Level,
		)
	}

	return storage.OilFillLevels, nil
}

func (s *Storage) GetOilFillLevel(id int64) (models.OilFillLevel, error) {
	levels, err := s.GetOilFillLevels()
	if err != nil {
		return models.OilFillLevel{}, err
	}

	for _, level := range levels {
		if level.Id == id {
			return level, nil
		}
	}

	return models.OilFillLevel{}, nil
}

func (s *Storage) UpdateOilFillLevel(oilFillLevel models.OilFillLevel) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	for i := range storage.OilFillLevels {
		if storage.OilFillLevels[i].Id == oilFillLevel.Id {
			storage.OilFillLevels[i] = oilFillLevel
			break
		}
	}

	return s.saveStorage(storage)
}

func (s *Storage) DeleteOilFillLevel(id int64) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	for i := range storage.OilFillLevels {
		if storage.OilFillLevels[i].Id == id {
			storage.OilFillLevels = slices.Delete(storage.OilFillLevels, i, i+1)
			break
		}
	}

	return s.saveStorage(storage)
}

// Helpers

func calculateOilFillPercentage(cm float64) float64 {
	if cm < 0 {
		cm = 0
	}

	percentage := (cm / models.OilTankHeightCm) * 100

	if percentage < 0 {
		return 0
	}
	if percentage > 100 {
		return 100
	}

	return percentage
}

func sortOil(oil []models.Oil) []models.Oil {
	sort.Slice(oil, func(i, j int) bool {
		return oil[j].Date < oil[i].Date
	})

	return oil
}

func sortOilFillLevels(oilFillLevels []models.OilFillLevel) []models.OilFillLevel {
	sort.Slice(oilFillLevels, func(i, j int) bool {
		return oilFillLevels[j].Date < oilFillLevels[i].Date
	})

	return oilFillLevels
}
