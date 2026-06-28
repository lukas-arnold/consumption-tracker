package storage

import (
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestAddOil(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	err := store.AddOil(
		models.OilInput{
			Date:   "2024-01-01",
			Volume: 100,
			Costs:  50,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	oils, err :=
		store.GetOil()

	if err != nil {
		t.Fatal(err)
	}

	if len(oils) != 1 {
		t.Fatal(
			"oil not added",
		)
	}

	if oils[0].Volume != 100 {
		t.Fatalf(
			"got %f",
			oils[0].Volume,
		)
	}
}

func TestGetOilEntry(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOil(
		models.OilInput{
			Volume: 50,
		},
	)

	oils, _ :=
		store.GetOil()

	oil, err :=
		store.GetOilEntry(
			oils[0].Id,
		)

	if err != nil {
		t.Fatal(err)
	}

	if oil.Volume != 50 {
		t.Fatalf(
			"got %f",
			oil.Volume,
		)
	}
}

func TestGetOilEntryNotFound(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	oil, err :=
		store.GetOilEntry(999)

	if err != nil {
		t.Fatal(err)
	}

	if oil.Id != 0 {
		t.Fatal(
			"expected empty oil",
		)
	}
}

func TestUpdateOil(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOil(
		models.OilInput{
			Volume: 50,
		},
	)

	oils, _ :=
		store.GetOil()

	oil :=
		oils[0]

	oil.Volume = 100

	err := store.UpdateOil(
		oil,
	)

	if err != nil {
		t.Fatal(err)
	}

	updated, _ :=
		store.GetOilEntry(
			oil.Id,
		)

	if updated.Volume != 100 {
		t.Fatalf(
			"got %f",
			updated.Volume,
		)
	}
}

func TestDeleteOil(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOil(
		models.OilInput{},
	)

	oils, _ :=
		store.GetOil()

	err := store.DeleteOil(
		oils[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	result, _ :=
		store.GetOil()

	if len(result) != 0 {
		t.Fatal(
			"oil not deleted",
		)
	}
}

func TestSortOil(t *testing.T) {

	oil := []models.Oil{
		{
			OilInput: models.OilInput{
				Date: "2024-01-01",
			},
		},
		{
			OilInput: models.OilInput{
				Date: "2025-01-01",
			},
		},
	}

	sorted :=
		sortOil(oil)

	if sorted[0].Date != "2025-01-01" {
		t.Fatal(
			"wrong order",
		)
	}
}

func TestAddOilFillLevel(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	err := store.AddOilFillLevel(
		models.OilFillLevelInput{
			Date:  "2024",
			Level: 75,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	levels, _ :=
		store.GetOilFillLevels()

	if len(levels) != 1 {
		t.Fatal(
			"level missing",
		)
	}
}

func TestGetOilFillLevelPercentage(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOilFillLevel(
		models.OilFillLevelInput{
			Level: 75,
		},
	)

	levels, err :=
		store.GetOilFillLevels()

	if err != nil {
		t.Fatal(err)
	}

	if levels[0].Percentage != 50 {
		t.Fatalf(
			"got %f",
			levels[0].Percentage,
		)
	}
}

func TestCalculateOilFillPercentage(t *testing.T) {

	if calculateOilFillPercentage(75) != 50 {
		t.Fatal(
			"wrong percentage",
		)
	}

	if calculateOilFillPercentage(-10) != 0 {
		t.Fatal(
			"expected minimum 0",
		)
	}

	if calculateOilFillPercentage(300) != 100 {
		t.Fatal(
			"expected maximum 100",
		)
	}
}

func TestDeleteOilFillLevel(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOilFillLevel(
		models.OilFillLevelInput{
			Level: 50,
		},
	)

	levels, _ :=
		store.GetOilFillLevels()

	err := store.DeleteOilFillLevel(
		levels[0].Id,
	)

	if err != nil {
		t.Fatal(err)
	}

	result, _ :=
		store.GetOilFillLevels()

	if len(result) != 0 {
		t.Fatal(
			"level not deleted",
		)
	}
}

func TestSortOilFillLevels(t *testing.T) {

	levels := []models.OilFillLevel{
		{
			OilFillLevelInput: models.OilFillLevelInput{
				Date: "2024",
			},
		},
		{
			OilFillLevelInput: models.OilFillLevelInput{
				Date: "2025",
			},
		},
	}

	sorted :=
		sortOilFillLevels(levels)

	if sorted[0].Date != "2025" {
		t.Fatal(
			"wrong order",
		)
	}
}
