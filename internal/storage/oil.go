package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func AddOil(oil models.OilInput) error {
	storage, err := GetUsageStorage()
	if err != nil {
		return err
	}
	newOil := models.Oil{Id: utils.Id(), OilInput: oil}
	storage.Oil = append(storage.Oil, newOil)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func GetOil() ([]models.Oil, error) {
	storage, err := GetUsageStorage()
	if err != nil {
		return nil, err
	}
	return storage.Oil, nil
}

func GetOilEntry(id int64) (models.Oil, error) {
	oils, err := GetOil()
	if err != nil {
		return models.Oil{}, err
	}
	var oil models.Oil
	for _, value := range oils {
		if value.Id == id {
			oil = value
		}
	}
	return oil, nil
}

func UpdateOil(oil models.Oil) error {
	storage, err := GetUsageStorage()
	if err != nil {
		return err
	}
	for i := range storage.Oil {
		if storage.Oil[i].Id == oil.Id {
			storage.Oil[i].Date = oil.Date
			storage.Oil[i].Volume = oil.Volume
			storage.Oil[i].Costs = oil.Costs
			storage.Oil[i].Retailer = oil.Retailer
			storage.Oil[i].Note = oil.Note
		}
	}
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOil(id int64) error {
	storage, err := GetUsageStorage()
	if err != nil {
		return err
	}
	var index int
	for i := range storage.Oil {
		if storage.Oil[i].Id == id {
			index = i
		}
	}
	storage.Oil = slices.Delete(storage.Oil, index, index+1)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func AddOilFillLevel(oilFillLevel models.OilFillLevelInput) error {
	storage, err := GetUsageStorage()
	if err != nil {
		return err
	}
	newOilFillLevel := models.OilFillLevel{Id: utils.Id(), OilFillLevelInput: oilFillLevel}
	storage.OilFillLevels = append(storage.OilFillLevels, newOilFillLevel)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func GetOilFillLevels() ([]models.OilFillLevel, error) {
	storage, err := GetUsageStorage()
	if err != nil {
		return nil, err
	}
	for i := range storage.OilFillLevels {
		storage.OilFillLevels[i].Percentage = calculateOilFillPercentage(storage.OilFillLevels[i].Level)
	}
	return storage.OilFillLevels, nil
}

func GetOilFillLevel(id int64) (models.OilFillLevel, error) {
	oilFillLevels, err := GetOilFillLevels()
	if err != nil {
		return models.OilFillLevel{}, err
	}
	var oilFillLevel models.OilFillLevel
	for _, value := range oilFillLevels {
		if value.Id == id {
			oilFillLevel = value
		}
	}
	return oilFillLevel, nil
}

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

func UpdateOilFillLevel(oilFillLevel models.OilFillLevel) error {
	storage, err := GetUsageStorage()
	if err != nil {
		return err
	}
	for i := range storage.OilFillLevels {
		if storage.OilFillLevels[i].Id == oilFillLevel.Id {
			storage.OilFillLevels[i].Date = oilFillLevel.Date
			storage.OilFillLevels[i].Level = oilFillLevel.Level
		}
	}
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOilFillLevel(id int64) error {
	storage, err := GetUsageStorage()
	if err != nil {
		return err
	}
	var index int
	for i := range storage.OilFillLevels {
		if storage.OilFillLevels[i].Id == id {
			index = i
		}
	}
	storage.OilFillLevels = slices.Delete(storage.OilFillLevels, index, index+1)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func sortOil(oil []models.Oil) []models.Oil {
	sort.Slice(oil, func(i, j int) bool {
		return oil[i].Date < oil[j].Date
	})
	return oil
}

func sortOilFillLevels(oilFillLevels []models.OilFillLevel) []models.OilFillLevel {
	sort.Slice(oilFillLevels, func(i, j int) bool {
		return oilFillLevels[i].Date < oilFillLevels[j].Date
	})
	return oilFillLevels
}
