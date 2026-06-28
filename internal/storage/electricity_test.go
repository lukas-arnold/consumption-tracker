package storage

import (
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestAddElectricity(t *testing.T) {
	store := testStorage(t)

	err := store.AddElectricity(models.ElectricityInput{
		TimeFrom:    "2024-01-01",
		TimeTo:      "2024-12-31",
		Consumption: 1000,
		Costs:       200,
	})
	if err != nil {
		t.Fatal(err)
	}

	electricities, err := store.GetElectricities()
	if err != nil {
		t.Fatal(err)
	}

	if len(electricities) != 1 {
		t.Fatal("electricity not added")
	}

	if electricities[0].Consumption != 1000 {
		t.Fatalf("got %f", electricities[0].Consumption)
	}
}

func TestGetElectricity(t *testing.T) {
	store := testStorage(t)

	_ = store.AddElectricity(models.ElectricityInput{
		Consumption: 500,
		Costs:       100,
	})

	electricities, _ := store.GetElectricities()
	electricity, err := store.GetElectricity(electricities[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	if electricity.Consumption != 500 {
		t.Fatalf("got %f", electricity.Consumption)
	}
}

func TestGetElectricityNotFound(t *testing.T) {
	store := testStorage(t)

	electricity, err := store.GetElectricity(999)
	if err != nil {
		t.Fatal(err)
	}

	if electricity.Id != 0 {
		t.Fatal("expected empty electricity")
	}
}

func TestUpdateElectricity(t *testing.T) {
	store := testStorage(t)

	_ = store.AddElectricity(models.ElectricityInput{Consumption: 500})

	electricities, _ := store.GetElectricities()
	electricity := electricities[0]
	electricity.Consumption = 1000

	err := store.UpdateElectricity(electricity)
	if err != nil {
		t.Fatal(err)
	}

	updated, _ := store.GetElectricity(electricity.Id)
	if updated.Consumption != 1000 {
		t.Fatalf("got %f", updated.Consumption)
	}
}

func TestDeleteElectricity(t *testing.T) {
	store := testStorage(t)

	_ = store.AddElectricity(models.ElectricityInput{Consumption: 100})

	electricities, _ := store.GetElectricities()
	err := store.DeleteElectricity(electricities[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := store.GetElectricities()
	if len(result) != 0 {
		t.Fatal("electricity not deleted")
	}
}

func TestSortElectricity(t *testing.T) {
	electricity := []models.Electricity{
		{ElectricityInput: models.ElectricityInput{TimeFrom: "2024-01-01"}},
		{ElectricityInput: models.ElectricityInput{TimeFrom: "2025-01-01"}},
	}

	sorted := sortElectricity(electricity)

	if sorted[0].TimeFrom != "2025-01-01" {
		t.Fatal("wrong order")
	}
}
