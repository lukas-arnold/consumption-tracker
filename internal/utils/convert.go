package utils

import (
	"encoding/json"
	"strconv"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func ConvertUsageStorageToBytes(storage models.UsageStorage) ([]byte, error) {
	bytes, err := json.Marshal(storage)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ConvertBytesToUsageStorage(bytes []byte) (models.UsageStorage, error) {
	var storage models.UsageStorage
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
