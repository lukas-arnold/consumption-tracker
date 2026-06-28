package storage

import (
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestGetElectricityCharts(t *testing.T) {
	store := testStorage(t)
	store.AddElectricity(models.ElectricityInput{
		TimeFrom:    "2024-01-01",
		TimeTo:      "2024-01-11",
		Consumption: 100,
		Costs:       50,
	})

	charts, err := store.GetElectricityCharts()
	if err != nil {
		t.Fatal(err)
	}

	if len(charts.Consumption.Labels) != 1 || charts.Consumption.Labels[0] != "2024" {
		t.Fatal("expected 2024 label")
	}
	if charts.Consumption.Sets[0].Data[0] != 100 {
		t.Fatalf("got consumption %f", charts.Consumption.Sets[0].Data[0])
	}
	if charts.Price.Sets[0].Data[0] != 0.5 {
		t.Fatalf("got price %f", charts.Price.Sets[0].Data[0])
	}
}

func TestGetElectricityChartsMultiYear(t *testing.T) {
	store := testStorage(t)
	store.AddElectricity(models.ElectricityInput{
		TimeFrom:    "2024-12-01",
		TimeTo:      "2025-01-01",
		Consumption: 310,
		Costs:       31,
	})

	charts, err := store.GetElectricityCharts()
	if err != nil {
		t.Fatal(err)
	}

	if len(charts.Consumption.Labels) != 1 || charts.Consumption.Labels[0] != "2024" {
		t.Fatalf("expected only 2024, got %#v", charts.Consumption.Labels)
	}
}

func TestGetElectricityChartsInvalidDate(t *testing.T) {
	store := testStorage(t)
	store.AddElectricity(models.ElectricityInput{TimeFrom: "broken", TimeTo: "2024-01-01"})

	if _, err := store.GetElectricityCharts(); err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestGetOilCharts(t *testing.T) {
	store := testStorage(t)
	store.AddOil(models.OilInput{Date: "2024-01-01", Volume: 100, Costs: 200})

	charts, err := store.GetOilCharts()
	if err != nil {
		t.Fatal(err)
	}

	if charts.Consumption.Labels[0] != "2024" || charts.Consumption.Sets[0].Data[0] != 100 {
		t.Fatal("unexpected oil consumption data")
	}
	if charts.Price.Sets[0].Data[0] != 2 {
		t.Fatalf("expected price 2, got %f", charts.Price.Sets[0].Data[0])
	}
}

func TestGetOilChartsFillLevel(t *testing.T) {
	store := testStorage(t)
	store.AddOilFillLevel(models.OilFillLevelInput{Date: "2024-01-01", Level: 75})

	charts, err := store.GetOilCharts()
	if err != nil {
		t.Fatal(err)
	}

	if len(charts.FillLevel.Labels) != 1 {
		t.Fatal("fill level missing")
	}
	if charts.FillLevel.Sets[1].Data[0] != 50 {
		t.Fatalf("got %f", charts.FillLevel.Sets[1].Data[0])
	}
}

func TestGetWaterCharts(t *testing.T) {
	store := testStorage(t)
	store.AddWater(models.WaterInput{
		Year: 2024, VolumeWater: 100, VolumeWastewater: 50, VolumeRainwater: 10,
		CostsWater: 200, CostsWastewater: 50, CostsRainwater: 10, FixedPrice: 100,
	})

	charts, err := store.GetWaterCharts()
	if err != nil {
		t.Fatal(err)
	}

	if charts.Consumption.Labels[0] != "2024" {
		t.Fatal("wrong year")
	}
	if charts.Consumption.Sets[0].Data[0] != 100 || charts.Consumption.Sets[1].Data[0] != 60 {
		t.Fatal("wrong volume calculations")
	}
	if charts.Price.Sets[0].Data[0] != 2 || charts.Price.Sets[3].Data[0] != 100 {
		t.Fatal("wrong price calculations")
	}
}
