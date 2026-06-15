package storage

import (
	"os"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func saveStorage(storage models.ConsumptionStorage) error {
	storage = sortStorage(storage)
	bytes, err := utils.ConvertConsumptionStorageToBytes(storage)
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
		saveStorage(models.ConsumptionStorage{})
	}
}

func getConsumptionStorage() (models.ConsumptionStorage, error) {
	bytes, err := readStorage()
	if err != nil {
		return models.ConsumptionStorage{}, err
	}
	storage, err := utils.ConvertBytesToConsumptionStorage(bytes)
	if err != nil {
		return models.ConsumptionStorage{}, err
	}
	return storage, nil
}

func sortStorage(storage models.ConsumptionStorage) models.ConsumptionStorage {
	storage.Electricity = sortElectricity(storage.Electricity)
	storage.Oil = sortOil(storage.Oil)
	storage.OilFillLevels = sortOilFillLevels(storage.OilFillLevels)
	storage.Water = sortWater(storage.Water)
	return storage
}
