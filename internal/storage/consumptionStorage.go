package storage

import (
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func GetConsumptionStorage() (models.ConsumptionStorage, error) {
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
