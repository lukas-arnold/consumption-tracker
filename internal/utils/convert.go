package utils

import (
	"encoding/json"
	"strconv"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func ConvertConsumptionStorageToBytes(storage models.ConsumptionStorage) ([]byte, error) {
	bytes, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ConvertBytesToConsumptionStorage(bytes []byte) (models.ConsumptionStorage, error) {
	var storage models.ConsumptionStorage
	err := json.Unmarshal(bytes, &storage)
	if err != nil {
		return storage, err
	}
	return storage, nil
}

func ConvertId(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ConvertFloat(fStr string) (float64, error) {
	f, err := strconv.ParseFloat(fStr, 64)
	if err != nil {
		return -1, err
	}
	return f, nil
}

func ConvertInt(iStr string) (int, error) {
	i, err := strconv.Atoi(iStr)
	if err != nil {
		return -1, err
	}
	return i, nil
}
