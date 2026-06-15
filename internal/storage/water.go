package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func AddWater(water models.WaterInput) error {
	storage, err := getConsumptionStorage()
	if err != nil {
		return err
	}
	newWater := models.Water{Id: utils.Id(), WaterInput: water}
	storage.Water = append(storage.Water, newWater)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func GetWater() ([]models.Water, error) {
	storage, err := getConsumptionStorage()
	if err != nil {
		return nil, err
	}
	return storage.Water, nil
}

func GetWaterEntry(id int64) (models.Water, error) {
	waters, err := GetWater()
	if err != nil {
		return models.Water{}, err
	}
	var water models.Water
	for _, value := range waters {
		if value.Id == id {
			water = value
		}
	}
	return water, nil
}

func UpdateWater(water models.Water) error {
	storage, err := getConsumptionStorage()
	if err != nil {
		return err
	}
	for i := range storage.Water {
		if storage.Water[i].Id == water.Id {
			storage.Water[i].Year = water.Year
			storage.Water[i].VolumeWater = water.VolumeWater
			storage.Water[i].VolumeWastewater = water.VolumeWastewater
			storage.Water[i].VolumeRainwater = water.VolumeRainwater
			storage.Water[i].CostsWater = water.CostsWater
			storage.Water[i].CostsWastewater = water.CostsWastewater
			storage.Water[i].CostsRainwater = water.CostsRainwater
			storage.Water[i].Payments = water.Payments
			storage.Water[i].FixedPrice = water.FixedPrice
			storage.Water[i].Note = water.Note
		}
	}
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteWater(id int64) error {
	storage, err := getConsumptionStorage()
	if err != nil {
		return err
	}
	var index int
	for i := range storage.Water {
		if storage.Water[i].Id == id {
			index = i
		}
	}
	storage.Water = slices.Delete(storage.Water, index, index+1)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func sortWater(water []models.Water) []models.Water {
	sort.Slice(water, func(i, j int) bool {
		return water[j].Year < water[i].Year
	})
	return water
}
