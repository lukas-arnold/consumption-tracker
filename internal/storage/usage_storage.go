package storage

import (
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func GetUsageStorage() (models.UsageStorage, error) {
	bytes, err := readStorage()
	if err != nil {
		return models.UsageStorage{}, err
	}
	storage, err := utils.ConvertBytesToUsageStorage(bytes)
	if err != nil {
		return models.UsageStorage{}, err
	}
	return storage, nil
}
