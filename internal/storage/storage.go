package storage

import (
	"os"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func saveStorage(storage models.UsageStorage) error {
	storage = sortStorage(storage)
	bytes, err := utils.ConvertUsageStorageToBytes(storage)
	if err != nil {
		return err
	}
	err = os.WriteFile(configs.GetStorageFile(), []byte(bytes), 0666)
	if err != nil {
		return err
	}
	return nil
}

func readStorage() ([]byte, error) {
	checkStorage()
	bytes, err := os.ReadFile(configs.GetStorageFile())
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func checkStorage() {
	_, err := os.ReadFile(configs.GetStorageFile())
	if err != nil {
		saveStorage(models.UsageStorage{})
	}
}

func sortStorage(storage models.UsageStorage) models.UsageStorage {
	storage.Electricity = sortElectricity(storage.Electricity)
	storage.Oil = sortOil(storage.Oil)
	storage.OilFillLevels = sortOilFillLevels(storage.OilFillLevels)
	storage.Water = sortWater(storage.Water)
	return storage
}
