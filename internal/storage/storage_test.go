package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	return New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)
}

func TestStorageCreatesFile(t *testing.T) {

	store := newTestStorage(t)

	err := store.checkStorage()

	if err != nil {
		t.Fatal(err)
	}

	data, err := store.readStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal("expected storage file")
	}
}

func TestStorageCreatesMissingDirectory(t *testing.T) {

	dir := t.TempDir()

	store := New(
		filepath.Join(
			dir,
			"nested",
			"storage.json",
		),
	)

	err := store.checkStorage()

	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(
		filepath.Join(
			dir,
			"nested",
		),
	)

	if err != nil {
		t.Fatal("expected directory to be created")
	}
}

func TestSaveStoragePersistsData(t *testing.T) {

	store := newTestStorage(t)

	input := models.ConsumptionStorage{
		Electricity: []models.Electricity{
			{
				ElectricityInput: models.ElectricityInput{
					TimeFrom:    "2026-01-01",
					TimeTo:      "2026-12-31",
					Consumption: 1234,
				},
			},
		},
	}

	err := store.saveStorage(input)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getConsumptionStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Electricity) != 1 {
		t.Fatalf(
			"expected 1 electricity entry got %d",
			len(result.Electricity),
		)
	}

	if result.Electricity[0].Consumption != 1234 {
		t.Fatalf(
			"wrong consumption %f",
			result.Electricity[0].Consumption,
		)
	}
}

func TestSaveStorageSortsElectricity(t *testing.T) {

	store := newTestStorage(t)

	err := store.saveStorage(
		models.ConsumptionStorage{
			Electricity: []models.Electricity{
				{
					ElectricityInput: models.ElectricityInput{
						TimeFrom: "2026-01-01",
					},
				},
				{
					ElectricityInput: models.ElectricityInput{
						TimeFrom: "2026-02-01",
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getConsumptionStorage()

	if err != nil {
		t.Fatal(err)
	}

	if result.Electricity[0].TimeFrom != "2026-02-01" {
		t.Fatalf(
			"expected sorted electricity got %s",
			result.Electricity[0].TimeFrom,
		)
	}
}

func TestSaveStorageSortsOil(t *testing.T) {

	store := newTestStorage(t)

	err := store.saveStorage(
		models.ConsumptionStorage{
			Oil: []models.Oil{
				{
					OilInput: models.OilInput{
						Date: "2026-01-01",
					},
				},
				{
					OilInput: models.OilInput{
						Date: "2026-02-01",
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getConsumptionStorage()

	if err != nil {
		t.Fatal(err)
	}

	if result.Oil[0].Date != "2026-02-01" {
		t.Fatalf(
			"expected sorted oil got %s",
			result.Oil[0].Date,
		)
	}
}

func TestSaveStorageSortsOilFillLevels(t *testing.T) {

	store := newTestStorage(t)

	err := store.saveStorage(
		models.ConsumptionStorage{
			OilFillLevels: []models.OilFillLevel{
				{
					OilFillLevelInput: models.OilFillLevelInput{
						Date: "2026-01-01",
					},
				},
				{
					OilFillLevelInput: models.OilFillLevelInput{
						Date: "2026-02-01",
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getConsumptionStorage()

	if err != nil {
		t.Fatal(err)
	}

	if result.OilFillLevels[0].Date != "2026-02-01" {
		t.Fatalf(
			"expected sorted oil fill levels got %s",
			result.OilFillLevels[0].Date,
		)
	}
}

func TestSaveStorageSortsWater(t *testing.T) {

	store := newTestStorage(t)

	err := store.saveStorage(
		models.ConsumptionStorage{
			Water: []models.Water{
				{
					WaterInput: models.WaterInput{
						Year: 2025,
					},
				},
				{
					WaterInput: models.WaterInput{
						Year: 2026,
					},
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	result, err := store.getConsumptionStorage()

	if err != nil {
		t.Fatal(err)
	}

	if result.Water[0].Year != 2026 {
		t.Fatalf(
			"expected sorted water got %d",
			result.Water[0].Year,
		)
	}
}

func TestReadStorageCreatesEmptyStorage(t *testing.T) {

	store := newTestStorage(t)

	data, err := store.readStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Fatal(
			"expected initialized storage",
		)
	}
}

func TestGetConsumptionStorageEmpty(t *testing.T) {

	store := newTestStorage(t)

	result, err := store.getConsumptionStorage()

	if err != nil {
		t.Fatal(err)
	}

	if len(result.Electricity) != 0 {
		t.Fatal("expected no electricity")
	}

	if len(result.Oil) != 0 {
		t.Fatal("expected no oil")
	}

	if len(result.OilFillLevels) != 0 {
		t.Fatal("expected no oil fill levels")
	}

	if len(result.Water) != 0 {
		t.Fatal("expected no water")
	}
}
