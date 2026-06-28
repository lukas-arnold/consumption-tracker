package utils

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func ConvertConsumptionStorageToBytes(storage models.ConsumptionStorage) ([]byte, error) {
	return json.Marshal(storage)
}

func ConvertBytesToConsumptionStorage(bytes []byte) (models.ConsumptionStorage, error) {
	var storage models.ConsumptionStorage
	err := json.Unmarshal(bytes, &storage)
	return storage, err
}

func ConvertId(idStr string) (int64, error) {
	return strconv.ParseInt(idStr, 10, 64)
}

func ConvertFloat(fStr string) (float64, error) {
	return strconv.ParseFloat(fStr, 64)
}

func ConvertInt(iStr string) (int, error) {
	return strconv.Atoi(iStr)
}

func ConvertTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02", timeStr)
}
