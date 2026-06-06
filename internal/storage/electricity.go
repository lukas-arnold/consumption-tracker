package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func AddElectricity(electricity models.ElectricityInput) error {
	storage, err := GetConsumptionStorage()
	if err != nil {
		return err
	}
	newElectricity := models.Electricity{Id: utils.Id(), ElectricityInput: electricity}
	storage.Electricity = append(storage.Electricity, newElectricity)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func GetElectricities() ([]models.Electricity, error) {
	storage, err := GetConsumptionStorage()
	if err != nil {
		return nil, err
	}
	return storage.Electricity, nil
}

func GetElectricity(id int64) (models.Electricity, error) {
	electricities, err := GetElectricities()
	if err != nil {
		return models.Electricity{}, err
	}
	var electricity models.Electricity
	for _, value := range electricities {
		if value.Id == id {
			electricity = value
		}
	}
	return electricity, nil
}

func UpdateElectricity(electricity models.Electricity) error {
	storage, err := GetConsumptionStorage()
	if err != nil {
		return err
	}
	for i := range storage.Electricity {
		if storage.Electricity[i].Id == electricity.Id {
			storage.Electricity[i].TimeFrom = electricity.TimeFrom
			storage.Electricity[i].TimeTo = electricity.TimeTo
			storage.Electricity[i].Consumption = electricity.Consumption
			storage.Electricity[i].Costs = electricity.Costs
			storage.Electricity[i].Retailer = electricity.Retailer
			storage.Electricity[i].Payments = electricity.Payments
			storage.Electricity[i].Note = electricity.Note
		}
	}
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteElectricity(id int64) error {
	storage, err := GetConsumptionStorage()
	if err != nil {
		return err
	}
	var index int
	for i := range storage.Electricity {
		if storage.Electricity[i].Id == id {
			index = i
		}
	}
	storage.Electricity = slices.Delete(storage.Electricity, index, index+1)
	err = saveStorage(storage)
	if err != nil {
		return err
	}
	return nil
}

func sortElectricity(electricity []models.Electricity) []models.Electricity {
	sort.Slice(electricity, func(i, j int) bool {
		return electricity[j].TimeFrom < electricity[i].TimeFrom
	})
	return electricity
}
