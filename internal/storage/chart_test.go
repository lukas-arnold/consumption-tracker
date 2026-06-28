package storage

import (
	"path/filepath"
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestGetElectricityCharts(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	err := store.AddElectricity(
		models.ElectricityInput{
			TimeFrom:    "2024-01-01",
			TimeTo:      "2024-01-11",
			Consumption: 100,
			Costs:       50,
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	charts, err :=
		store.GetElectricityCharts()

	if err != nil {
		t.Fatal(err)
	}

	if len(charts.Consumption.Labels) != 1 {
		t.Fatal(
			"expected one year",
		)
	}

	if charts.Consumption.Labels[0] != "2024" {
		t.Fatal(
			"wrong year",
		)
	}

	consumption :=
		charts.Consumption.Sets[0].Data[0]

	if consumption != 100 {
		t.Fatalf(
			"got %f",
			consumption,
		)
	}

	price :=
		charts.Price.Sets[0].Data[0]

	if price != 0.5 {
		t.Fatalf(
			"got %f",
			price,
		)
	}
}

func TestGetElectricityChartsMultiYear(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddElectricity(
		models.ElectricityInput{
			TimeFrom:    "2024-12-01",
			TimeTo:      "2025-01-01",
			Consumption: 310,
			Costs:       31,
		},
	)

	charts, err :=
		store.GetElectricityCharts()

	if err != nil {
		t.Fatal(err)
	}

	if len(charts.Consumption.Labels) != 2 {
		t.Fatal(
			"expected split years",
		)
	}

	if charts.Consumption.Labels[0] != "2024" {
		t.Fatal(
			"wrong order",
		)
	}
}

func TestGetElectricityChartsInvalidDate(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddElectricity(
		models.ElectricityInput{
			TimeFrom: "broken",
			TimeTo:   "2024-01-01",
		},
	)

	_, err :=
		store.GetElectricityCharts()

	if err == nil {
		t.Fatal(
			"expected date error",
		)
	}
}

func TestGetOilCharts(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOil(
		models.OilInput{
			Date:   "2024-01-01",
			Volume: 100,
			Costs:  200,
		},
	)

	charts, err :=
		store.GetOilCharts()

	if err != nil {
		t.Fatal(err)
	}

	if charts.Consumption.Labels[0] != "2024" {
		t.Fatal(
			"wrong year",
		)
	}

	if charts.Consumption.Sets[0].Data[0] != 100 {
		t.Fatalf(
			"got %f",
			charts.Consumption.Sets[0].Data[0],
		)
	}

	if charts.Price.Sets[0].Data[0] != 2 {
		t.Fatalf(
			"got %f",
			charts.Price.Sets[0].Data[0],
		)
	}
}

func TestGetOilChartsFillLevel(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddOilFillLevel(
		models.OilFillLevelInput{
			Date:  "2024-01-01",
			Level: 75,
		},
	)

	charts, err :=
		store.GetOilCharts()

	if err != nil {
		t.Fatal(err)
	}

	if len(charts.FillLevel.Labels) != 1 {
		t.Fatal(
			"fill level missing",
		)
	}

	if charts.FillLevel.Sets[1].Data[0] != 50 {
		t.Fatalf(
			"got %f",
			charts.FillLevel.Sets[1].Data[0],
		)
	}
}

func TestGetWaterCharts(t *testing.T) {

	store := New(
		filepath.Join(
			t.TempDir(),
			"storage.json",
		),
	)

	store.AddWater(
		models.WaterInput{
			Year:             2024,
			VolumeWater:      100,
			VolumeWastewater: 50,
			VolumeRainwater:  10,
			CostsWater:       200,
			CostsWastewater:  50,
			CostsRainwater:   10,
			FixedPrice:       100,
		},
	)

	charts, err :=
		store.GetWaterCharts()

	if err != nil {
		t.Fatal(err)
	}

	if charts.Consumption.Labels[0] != "2024" {
		t.Fatal(
			"wrong year",
		)
	}

	if charts.Consumption.Sets[0].Data[0] != 100 {
		t.Fatal(
			"wrong water volume",
		)
	}

	// wastewater + rainwater
	if charts.Consumption.Sets[1].Data[0] != 60 {
		t.Fatal(
			"wrong wastewater calculation",
		)
	}

	// water price
	if charts.Price.Sets[0].Data[0] != 2 {
		t.Fatalf(
			"got %f",
			charts.Price.Sets[0].Data[0],
		)
	}

	// fixed price average
	if charts.Price.Sets[3].Data[0] != 100 {
		t.Fatalf(
			"got %f",
			charts.Price.Sets[3].Data[0],
		)
	}
}
