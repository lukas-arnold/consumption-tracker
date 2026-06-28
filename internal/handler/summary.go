package handler

import (
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type ConsumptionSummary struct {
	TotalConsumption          float64
	TotalCosts                float64
	AverageConsumptionPerYear float64
	AverageCostPerUnit        float64
	YearsCount                int
}

func buildElectricitySummary(
	electricityEntries []models.Electricity,
) ConsumptionSummary {

	totalConsumption := 0.0
	totalCosts := 0.0

	for _, entry := range electricityEntries {
		totalConsumption += entry.Consumption
		totalCosts += entry.Costs
	}

	yearsCount := 0
	averageConsumption := 0.0

	if len(electricityEntries) > 0 {

		newest, err1 :=
			utils.ConvertTime(
				electricityEntries[0].TimeTo,
			)

		oldest, err2 :=
			utils.ConvertTime(
				electricityEntries[len(electricityEntries)-1].TimeFrom,
			)

		if err1 == nil && err2 == nil {

			yearsCount =
				newest.Year() -
					oldest.Year() +
					1

			years :=
				newest.Sub(oldest).
					Hours() /
					24 /
					365.25

			if years > 0 {
				averageConsumption =
					totalConsumption / years
			}
		}
	}

	averageCost := 0.0

	if totalConsumption > 0 {
		averageCost =
			totalCosts /
				totalConsumption
	}

	return ConsumptionSummary{
		TotalConsumption:          totalConsumption,
		TotalCosts:                totalCosts,
		AverageConsumptionPerYear: averageConsumption,
		AverageCostPerUnit:        averageCost,
		YearsCount:                yearsCount,
	}
}

func buildOilSummary(
	oilEntries []models.Oil,
) ConsumptionSummary {

	totalVolume := 0.0
	totalCosts := 0.0

	for _, entry := range oilEntries {
		totalVolume += entry.Volume
		totalCosts += entry.Costs
	}

	yearsCount := 0
	averageVolume := 0.0

	if len(oilEntries) > 0 {
		newest, err1 :=
			utils.ConvertTime(
				oilEntries[0].Date,
			)

		oldest, err2 :=
			utils.ConvertTime(
				oilEntries[len(oilEntries)-1].Date,
			)

		if err1 == nil && err2 == nil {
			yearsCount =
				newest.Year() -
					oldest.Year() +
					1

			years :=
				newest.Sub(oldest).
					Hours() /
					24 /
					365.25

			if years > 0 {
				averageVolume =
					totalVolume / years
			}
		}
	}

	averageCost := 0.0

	if totalVolume > 0 {
		averageCost =
			totalCosts /
				totalVolume
	}

	return ConsumptionSummary{
		TotalConsumption:          totalVolume,
		TotalCosts:                totalCosts,
		AverageConsumptionPerYear: averageVolume,
		AverageCostPerUnit:        averageCost,
		YearsCount:                yearsCount,
	}
}

func buildWaterSummary(
	waterEntries []models.Water,
) ConsumptionSummary {

	totalVolume := 0.0
	totalCosts := 0.0

	for _, entry := range waterEntries {

		totalVolume += entry.VolumeWater

		totalCosts +=
			entry.CostsWater +
				entry.CostsWastewater +
				entry.CostsRainwater +
				entry.FixedPrice
	}

	yearsCount := 0
	averageVolume := 0.0

	if len(waterEntries) > 0 {

		newestYear :=
			waterEntries[0].Year

		oldestYear :=
			waterEntries[len(waterEntries)-1].Year

		yearsCount =
			newestYear -
				oldestYear +
				1

		if yearsCount > 0 {
			averageVolume =
				totalVolume /
					float64(yearsCount)
		}
	}

	averageCost := 0.0

	if totalVolume > 0 {
		averageCost =
			totalCosts /
				totalVolume
	}

	return ConsumptionSummary{
		TotalConsumption:          totalVolume,
		TotalCosts:                totalCosts,
		AverageConsumptionPerYear: averageVolume,
		AverageCostPerUnit:        averageCost,
		YearsCount:                yearsCount,
	}
}
