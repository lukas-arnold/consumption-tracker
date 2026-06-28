package storage

import (
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestAddWater(t *testing.T) {
	store := testStorage(t)

	err := store.AddWater(models.WaterInput{
		Year:        2024,
		VolumeWater: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	waters, err := store.GetWater()
	if err != nil {
		t.Fatal(err)
	}

	if len(waters) != 1 {
		t.Fatal("water not added")
	}

	if waters[0].VolumeWater != 100 {
		t.Fatalf("got %f", waters[0].VolumeWater)
	}
}

func TestGetWater(t *testing.T) {
	store := testStorage(t)

	_ = store.AddWater(models.WaterInput{
		Year:        2024,
		VolumeWater: 250,
	})

	waters, err := store.GetWater()
	if err != nil {
		t.Fatal(err)
	}

	if waters[0].VolumeWater != 250 {
		t.Fatalf("got %f", waters[0].VolumeWater)
	}
}

func TestGetWaterEntry(t *testing.T) {
	store := testStorage(t)

	_ = store.AddWater(models.WaterInput{
		Year:        2024,
		VolumeWater: 50,
	})

	waters, _ := store.GetWater()
	water, err := store.GetWaterEntry(waters[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	if water.VolumeWater != 50 {
		t.Fatalf("got %f", water.VolumeWater)
	}
}

func TestGetWaterEntryNotFound(t *testing.T) {
	store := testStorage(t)

	water, err := store.GetWaterEntry(999)
	if err != nil {
		t.Fatal(err)
	}

	if water.Id != 0 {
		t.Fatal("expected empty water")
	}
}

func TestUpdateWater(t *testing.T) {
	store := testStorage(t)

	_ = store.AddWater(models.WaterInput{
		Year:        2024,
		VolumeWater: 50,
	})

	waters, _ := store.GetWater()
	water := waters[0]
	water.VolumeWater = 100

	err := store.UpdateWater(water)
	if err != nil {
		t.Fatal(err)
	}

	updated, _ := store.GetWaterEntry(water.Id)
	if updated.VolumeWater != 100 {
		t.Fatalf("got %f", updated.VolumeWater)
	}
}

func TestDeleteWater(t *testing.T) {
	store := testStorage(t)

	_ = store.AddWater(models.WaterInput{Year: 2024})

	waters, _ := store.GetWater()
	err := store.DeleteWater(waters[0].Id)
	if err != nil {
		t.Fatal(err)
	}

	result, _ := store.GetWater()
	if len(result) != 0 {
		t.Fatal("water not deleted")
	}
}

func TestSortWater(t *testing.T) {
	waters := []models.Water{
		{WaterInput: models.WaterInput{Year: 2023}},
		{WaterInput: models.WaterInput{Year: 2024}},
	}

	sorted := sortWater(waters)

	if sorted[0].Year != 2024 {
		t.Fatal("wrong order")
	}
}
