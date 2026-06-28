package storage

import (
	"slices"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func (s *Storage) AddElectricity(
	electricity models.ElectricityInput,
) error {

	storage, err := s.getConsumptionStorage()

	if err != nil {
		return err
	}

	newElectricity := models.Electricity{
		Id:               utils.Id(),
		ElectricityInput: electricity,
	}

	storage.Electricity = append(
		storage.Electricity,
		newElectricity,
	)

	return s.saveStorage(storage)
}

func (s *Storage) GetElectricities() ([]models.Electricity, error) {

	storage, err := s.getConsumptionStorage()

	if err != nil {
		return nil, err
	}

	return storage.Electricity, nil
}

func (s *Storage) GetElectricity(
	id int64,
) (models.Electricity, error) {

	electricities, err := s.GetElectricities()

	if err != nil {
		return models.Electricity{}, err
	}

	for _, electricity := range electricities {

		if electricity.Id == id {
			return electricity, nil
		}
	}

	return models.Electricity{}, nil
}

func (s *Storage) UpdateElectricity(
	electricity models.Electricity,
) error {

	storage, err := s.getConsumptionStorage()

	if err != nil {
		return err
	}

	for i := range storage.Electricity {

		if storage.Electricity[i].Id == electricity.Id {

			storage.Electricity[i] = electricity

			break
		}
	}

	return s.saveStorage(storage)
}

func (s *Storage) DeleteElectricity(
	id int64,
) error {

	storage, err := s.getConsumptionStorage()

	if err != nil {
		return err
	}

	for i := range storage.Electricity {

		if storage.Electricity[i].Id == id {

			storage.Electricity = slices.Delete(
				storage.Electricity,
				i,
				i+1,
			)

			break
		}
	}

	return s.saveStorage(storage)
}

func sortElectricity(
	electricity []models.Electricity,
) []models.Electricity {

	sort.Slice(
		electricity,
		func(i, j int) bool {
			return electricity[j].TimeFrom < electricity[i].TimeFrom
		},
	)

	return electricity
}
