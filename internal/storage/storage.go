package storage

import (
	"os"
	"path/filepath"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type Storage struct {
	file string
}

func New(file string) *Storage {
	return &Storage{
		file: file,
	}
}

func (s *Storage) saveStorage(storage models.ConsumptionStorage) error {
	storage = sortStorage(storage)

	bytes, err := utils.ConvertConsumptionStorageToBytes(storage)
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, []byte(bytes), 0666)
}

func (s *Storage) readStorage() ([]byte, error) {
	if err := s.checkStorage(); err != nil {
		return nil, err
	}

	return os.ReadFile(s.file)
}

func (s *Storage) checkStorage() error {
	_, err := os.ReadFile(s.file)
	if err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(s.file), 0755); err != nil {
		return err
	}

	return s.saveStorage(models.ConsumptionStorage{})
}

func (s *Storage) getConsumptionStorage() (models.ConsumptionStorage, error) {
	bytes, err := s.readStorage()
	if err != nil {
		return models.ConsumptionStorage{}, err
	}

	return utils.ConvertBytesToConsumptionStorage(bytes)
}

func sortStorage(storage models.ConsumptionStorage) models.ConsumptionStorage {
	storage.Electricity = sortElectricity(storage.Electricity)
	storage.Oil = sortOil(storage.Oil)
	storage.OilFillLevels = sortOilFillLevels(storage.OilFillLevels)
	storage.Water = sortWater(storage.Water)

	return storage
}
