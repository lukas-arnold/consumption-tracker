package storage

import (
	"fmt"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func GetElectricityCharts() (models.ElectricityCharts, error) {
	electricities, err := GetElectricities()
	if err != nil {
		return models.ElectricityCharts{}, err
	}

	type totals struct {
		consumption float64
		costs       float64
	}

	yearTotals := map[string]totals{}

	for _, v := range electricities {
		year := v.TimeFrom
		if len(year) >= 4 {
			year = year[:4]
		}

		t := yearTotals[year]
		t.consumption += v.Consumption
		t.costs += v.Costs
		yearTotals[year] = t
	}

	var labels []string
	for y := range yearTotals {
		labels = append(labels, y)
	}
	sort.Strings(labels)

	var consumption, costs, prices []float64

	for _, y := range labels {
		t := yearTotals[y]

		consumption = append(consumption, t.consumption)
		costs = append(costs, t.costs)

		price := 0.0
		if t.consumption > 0 {
			price = t.costs / t.consumption
		}
		prices = append(prices, price)
	}

	return models.ElectricityCharts{
		Consumption: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: "Consumption",
					Data:  consumption,
					Unit:  "kWh",
				},
				{
					Label: "Costs",
					Data:  costs,
					Unit:  "€",
				},
			},
		},
		Price: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: "Price per unit",
					Data:  prices,
					Unit:  "€/kWh",
				},
			},
		},
	}, nil
}

func GetOilCharts() (models.OilCharts, error) {
	oils, err := GetOil()
	if err != nil {
		return models.OilCharts{}, err
	}

	type totals struct {
		volume float64
		costs  float64
	}

	yearTotals := map[string]totals{}

	for _, v := range oils {
		year := v.Date
		if len(year) >= 4 {
			year = year[:4]
		}

		t := yearTotals[year]
		t.volume += v.Volume
		t.costs += v.Costs
		yearTotals[year] = t
	}

	var labels []string
	for y := range yearTotals {
		labels = append(labels, y)
	}
	sort.Strings(labels)

	var volumes, costs, prices []float64

	for _, y := range labels {
		t := yearTotals[y]

		volumes = append(volumes, t.volume)
		costs = append(costs, t.costs)

		price := 0.0
		if t.volume > 0 {
			price = t.costs / t.volume
		}
		prices = append(prices, price)
	}

	// Fill levels
	fillLevels, err := GetOilFillLevels()
	if err != nil {
		return models.OilCharts{}, err
	}

	sort.Slice(fillLevels, func(i, j int) bool {
		return fillLevels[i].Date < fillLevels[j].Date
	})

	var fillDates []string
	var fillValues []float64

	for _, v := range fillLevels {
		fillDates = append(fillDates, v.Date)
		fillValues = append(fillValues, v.Level)
	}

	return models.OilCharts{
		Consumption: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: "Consumption",
					Data:  volumes,
					Unit:  "l",
				},
				{
					Label: "Costs",
					Data:  costs,
					Unit:  "€",
				},
			},
		},
		Price: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: "Price per liter",
					Data:  prices,
					Unit:  "€/l",
				},
			},
		},
		FillLevel: models.ChartModel{
			Labels: fillDates,
			Sets: []models.ChartDataset{
				{
					Label: "Fill level",
					Data:  fillValues,
					Unit:  "cm",
				},
			},
		},
	}, nil
}

func GetWaterCharts() (models.WaterCharts, error) {
	waters, err := GetWater()
	if err != nil {
		return models.WaterCharts{}, err
	}

	type totals struct {
		volume float64
		costs  float64
	}

	yearTotals := map[int]totals{}

	for _, v := range waters {
		t := yearTotals[v.Year]
		t.volume += v.VolumeWater + v.VolumeWastewater + v.VolumeRainwater
		t.costs += v.CostsWater + v.CostsWastewater + v.CostsRainwater
		yearTotals[v.Year] = t
	}

	var years []int
	for y := range yearTotals {
		years = append(years, y)
	}
	sort.Ints(years)

	var labels []string
	var volumes, costs, prices []float64

	for _, y := range years {
		t := yearTotals[y]

		labels = append(labels, fmt.Sprintf("%d", y))
		volumes = append(volumes, t.volume)
		costs = append(costs, t.costs)

		price := 0.0
		if t.volume > 0 {
			price = t.costs / t.volume
		}
		prices = append(prices, price)
	}

	return models.WaterCharts{
		Consumption: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: "Consumption",
					Data:  volumes,
					Unit:  "m³",
				},
				{
					Label: "Costs",
					Data:  costs,
					Unit:  "€",
				},
			},
		},
		Price: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: "Price per m³",
					Data:  prices,
					Unit:  "€/m³",
				},
			},
		},
	}, nil
}
