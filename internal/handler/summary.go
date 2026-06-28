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

func buildElectricitySummary(entries []models.Electricity) ConsumptionSummary {
	totalConsumption := 0.0
	totalCosts := 0.0

	for _, e := range entries {
		totalConsumption += e.Consumption
		totalCosts += e.Costs
	}

	yearsCount := 0
	averageConsumption := 0.0

	if len(entries) > 0 {
		newest, err1 := utils.ConvertTime(entries[0].TimeTo)
		oldest, err2 := utils.ConvertTime(entries[len(entries)-1].TimeFrom)

		if err1 == nil && err2 == nil {
			yearsCount = newest.Year() - oldest.Year() + 1
			years := newest.Sub(oldest).Hours() / 24 / 365.25
			if years > 0 {
				averageConsumption = totalConsumption / years
			}
		}
	}

	averageCost := 0.0
	if totalConsumption > 0 {
		averageCost = totalCosts / totalConsumption
	}

	return ConsumptionSummary{
		TotalConsumption:          totalConsumption,
		TotalCosts:                totalCosts,
		AverageConsumptionPerYear: averageConsumption,
		AverageCostPerUnit:        averageCost,
		YearsCount:                yearsCount,
	}
}

func buildOilSummary(entries []models.Oil) ConsumptionSummary {
	totalVolume := 0.0
	totalCosts := 0.0

	for _, e := range entries {
		totalVolume += e.Volume
		totalCosts += e.Costs
	}

	yearsCount := 0
	averageVolume := 0.0

	if len(entries) > 0 {
		newest, err1 := utils.ConvertTime(entries[0].Date)
		oldest, err2 := utils.ConvertTime(entries[len(entries)-1].Date)

		if err1 == nil && err2 == nil {
			yearsCount = newest.Year() - oldest.Year() + 1
			years := newest.Sub(oldest).Hours() / 24 / 365.25
			if years > 0 {
				averageVolume = totalVolume / years
			}
		}
	}

	averageCost := 0.0
	if totalVolume > 0 {
		averageCost = totalCosts / totalVolume
	}

	return ConsumptionSummary{
		TotalConsumption:          totalVolume,
		TotalCosts:                totalCosts,
		AverageConsumptionPerYear: averageVolume,
		AverageCostPerUnit:        averageCost,
		YearsCount:                yearsCount,
	}
}

func buildWaterSummary(entries []models.Water) ConsumptionSummary {
	totalVolume := 0.0
	totalCosts := 0.0

	for _, e := range entries {
		totalVolume += e.VolumeWater
		totalCosts += e.CostsWater + e.CostsWastewater + e.CostsRainwater + e.FixedPrice
	}

	yearsCount := 0
	averageVolume := 0.0

	if len(entries) > 0 {
		newestYear := entries[0].Year
		oldestYear := entries[len(entries)-1].Year

		yearsCount = newestYear - oldestYear + 1
		if yearsCount > 0 {
			averageVolume = totalVolume / float64(yearsCount)
		}
	}

	averageCost := 0.0
	if totalVolume > 0 {
		averageCost = totalCosts / totalVolume
	}

	return ConsumptionSummary{
		TotalConsumption:          totalVolume,
		TotalCosts:                totalCosts,
		AverageConsumptionPerYear: averageVolume,
		AverageCostPerUnit:        averageCost,
		YearsCount:                yearsCount,
	}
}
