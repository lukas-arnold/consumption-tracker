package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func (s *Storage) AddWater(water models.WaterInput) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	newWater := models.Water{
		Id:         utils.Id(),
		WaterInput: water,
	}

	storage.Water = append(storage.Water, newWater)

	return s.saveStorage(storage)
}

func (s *Storage) GetWater() ([]models.Water, error) {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return nil, err
	}

	return storage.Water, nil
}

func (s *Storage) GetWaterEntry(id int64) (models.Water, error) {
	waters, err := s.GetWater()
	if err != nil {
		return models.Water{}, err
	}

	for _, water := range waters {
		if water.Id == id {
			return water, nil
		}
	}

	return models.Water{}, nil
}

func (s *Storage) UpdateWater(water models.Water) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	for i := range storage.Water {
		if storage.Water[i].Id == water.Id {
			storage.Water[i] = water
			break
		}
	}

	return s.saveStorage(storage)
}

func (s *Storage) DeleteWater(id int64) error {
	storage, err := s.getConsumptionStorage()
	if err != nil {
		return err
	}

	for i := range storage.Water {
		if storage.Water[i].Id == id {
			storage.Water = slices.Delete(storage.Water, i, i+1)
			break
		}
	}

	return s.saveStorage(storage)
}

func sortWater(water []models.Water) []models.Water {
	sort.Slice(water, func(i, j int) bool {
		return water[j].Year < water[i].Year
	})

	return water
}
