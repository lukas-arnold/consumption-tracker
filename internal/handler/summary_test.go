package handler

import (
	"testing"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func TestBuildElectricitySummary(t *testing.T) {

	entries := []models.Electricity{
		{
			ElectricityInput: models.ElectricityInput{
				TimeFrom:    "2025-01-01",
				TimeTo:      "2026-01-01",
				Consumption: 2000,
				Costs:       700,
			},
		},
		{
			ElectricityInput: models.ElectricityInput{
				TimeFrom:    "2024-01-01",
				TimeTo:      "2025-01-01",
				Consumption: 1000,
				Costs:       300,
			},
		},
	}

	result :=
		buildElectricitySummary(
			entries,
		)

	if result.TotalConsumption != 3000 {
		t.Fatalf(
			"got %f",
			result.TotalConsumption,
		)
	}

	if result.TotalCosts != 1000 {
		t.Fatalf(
			"got %f",
			result.TotalCosts,
		)
	}

	if result.AverageCostPerUnit != 1000.0/3000.0 {
		t.Fatalf(
			"got %f",
			result.AverageCostPerUnit,
		)
	}

	if result.YearsCount != 3 {
		t.Fatalf(
			"got %d",
			result.YearsCount,
		)
	}
}

func TestBuildOilSummary(t *testing.T) {

	entries := []models.Oil{
		{
			OilInput: models.OilInput{
				Date:   "2025-01-01",
				Volume: 300,
				Costs:  600,
			},
		},
		{
			OilInput: models.OilInput{
				Date:   "2024-01-01",
				Volume: 100,
				Costs:  200,
			},
		},
	}

	result :=
		buildOilSummary(
			entries,
		)

	if result.TotalConsumption != 400 {
		t.Fatalf(
			"got %f",
			result.TotalConsumption,
		)
	}

	if result.TotalCosts != 800 {
		t.Fatalf(
			"got %f",
			result.TotalCosts,
		)
	}

	if result.AverageCostPerUnit != 2 {
		t.Fatalf(
			"got %f",
			result.AverageCostPerUnit,
		)
	}

	if result.YearsCount != 2 {
		t.Fatalf(
			"got %d",
			result.YearsCount,
		)
	}
}

func TestBuildWaterSummary(t *testing.T) {

	entries := []models.Water{
		{
			WaterInput: models.WaterInput{
				Year:            2025,
				VolumeWater:     200,
				CostsWater:      100,
				CostsWastewater: 40,
				CostsRainwater:  20,
				FixedPrice:      10,
			},
		},
		{
			WaterInput: models.WaterInput{
				Year:            2024,
				VolumeWater:     100,
				CostsWater:      50,
				CostsWastewater: 20,
				CostsRainwater:  10,
				FixedPrice:      5,
			},
		},
	}

	result :=
		buildWaterSummary(
			entries,
		)

	if result.TotalConsumption != 300 {
		t.Fatalf(
			"got %f",
			result.TotalConsumption,
		)
	}

	if result.TotalCosts != 255 {
		t.Fatalf(
			"got %f",
			result.TotalCosts,
		)
	}

	if result.AverageCostPerUnit != 0.85 {
		t.Fatalf(
			"got %f",
			result.AverageCostPerUnit,
		)
	}

	if result.YearsCount != 2 {
		t.Fatalf(
			"got %d",
			result.YearsCount,
		)
	}

	if result.AverageConsumptionPerYear != 150 {
		t.Fatalf(
			"got %f",
			result.AverageConsumptionPerYear,
		)
	}
}

func TestBuildSummaryEmpty(t *testing.T) {

	tests := []ConsumptionSummary{
		buildElectricitySummary(nil),
		buildOilSummary(nil),
		buildWaterSummary(nil),
	}

	for _, result := range tests {

		if result.TotalConsumption != 0 {
			t.Fatal(
				"expected empty consumption",
			)
		}

		if result.TotalCosts != 0 {
			t.Fatal(
				"expected empty costs",
			)
		}

		if result.YearsCount != 0 {
			t.Fatal(
				"expected empty years",
			)
		}
	}
}
