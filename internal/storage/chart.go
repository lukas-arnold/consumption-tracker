package storage

import (
	"fmt"
	"sort"
	"time"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
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

		from, err := time.Parse("2006-01-02", v.TimeFrom)
		if err != nil {
			return models.ElectricityCharts{}, err
		}

		to, err := time.Parse("2006-01-02", v.TimeTo)
		if err != nil {
			return models.ElectricityCharts{}, err
		}

		// Number of days in the period
		totalDays := to.Sub(from).Hours() / 24
		if totalDays <= 0 {
			return models.ElectricityCharts{}, err
		}

		// Period can span several years
		for year := from.Year(); year <= to.Year(); year++ {

			yearStart := time.Date(
				year, 1, 1,
				0, 0, 0, 0,
				time.Local,
			)

			yearEnd := time.Date(
				year+1, 1, 1,
				0, 0, 0, 0,
				time.Local,
			)

			partStart := from
			if partStart.Before(yearStart) {
				partStart = yearStart
			}

			partEnd := to
			if partEnd.After(yearEnd) {
				partEnd = yearEnd
			}

			daysInYear := partEnd.Sub(partStart).Hours() / 24

			if daysInYear <= 0 {
				return models.ElectricityCharts{}, err
			}

			// This year's share of the period
			ratio := daysInYear / totalDays

			key := fmt.Sprintf("%d", year)

			t := yearTotals[key]

			t.consumption += v.Consumption * ratio
			t.costs += v.Costs * ratio

			yearTotals[key] = t
		}
	}

	var labels []string

	for y := range yearTotals {
		labels = append(labels, y)
	}

	sort.Strings(labels)

	var consumption []float64
	var costs []float64
	var prices []float64

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
					Label: language.T(configs.GetLanguage(), "consumption"),
					Data:  consumption,
					Unit:  "kWh",
					YAxis: "y",
				},
				{
					Label: language.T(configs.GetLanguage(), "costs"),
					Data:  costs,
					Unit:  "€",
					YAxis: "y1",
				},
			},
		},

		Price: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: language.T(configs.GetLanguage(), "pricePerUnit"),
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
	var fillPercentages []float64

	for _, v := range fillLevels {
		fillDates = append(fillDates, v.Date)
		fillValues = append(fillValues, v.Level)
		fillPercentages = append(fillPercentages, v.Percentage)
	}

	return models.OilCharts{
		Consumption: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: language.T(configs.GetLanguage(), "volume"),
					Data:  volumes,
					Unit:  "l",
					YAxis: "y",
				},
				{
					Label: language.T(configs.GetLanguage(), "costs"),
					Data:  costs,
					Unit:  "€",
					YAxis: "y1",
				},
			},
		},
		Price: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: language.T(configs.GetLanguage(), "pricePerUnit"),
					Data:  prices,
					Unit:  "€/l",
				},
			},
		},
		FillLevel: models.ChartModel{
			Labels: fillDates,
			Sets: []models.ChartDataset{
				{
					Label: language.T(configs.GetLanguage(), "levelCm"),
					Data:  fillValues,
					Unit:  "cm",
					YAxis: "y",
				},
				{
					Label: language.T(configs.GetLanguage(), "level%"),
					Data:  fillPercentages,
					Unit:  "%",
					YAxis: "y1",
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
		volumeWater      float64
		volumeWastewater float64
		volumeRainwater  float64

		costsWater      float64
		costsWastewater float64
		costsRainwater  float64

		fixedPriceSum float64
		fixedCount    int
	}

	yearTotals := map[int]totals{}

	for _, v := range waters {
		t := yearTotals[v.Year]

		t.volumeWater += v.VolumeWater
		t.volumeWastewater += v.VolumeWastewater
		t.volumeRainwater += v.VolumeRainwater

		t.costsWater += v.CostsWater
		t.costsWastewater += v.CostsWastewater
		t.costsRainwater += v.CostsRainwater

		if v.FixedPrice > 0 {
			t.fixedPriceSum += v.FixedPrice
			t.fixedCount++
		}

		yearTotals[v.Year] = t
	}

	var years []int
	for y := range yearTotals {
		years = append(years, y)
	}
	sort.Ints(years)

	var labels []string

	// consumption chart
	var volWater, volWaste, costWater, costWaste []float64

	// price chart
	var priceWater, priceWaste, priceRain, fixedPrices []float64

	for _, y := range years {
		t := yearTotals[y]
		labels = append(labels, fmt.Sprintf("%d", y))

		// Wastewater = wastewater + rainwater
		wasteVolume := t.volumeWastewater + t.volumeRainwater
		wasteCosts := t.costsWastewater + t.costsRainwater

		volWater = append(volWater, t.volumeWater)
		volWaste = append(volWaste, wasteVolume)

		costWater = append(costWater, t.costsWater)
		costWaste = append(costWaste, wasteCosts)

		// --- Prices ---
		pw := 0.0
		if t.volumeWater > 0 {
			pw = t.costsWater / t.volumeWater
		}

		pww := 0.0
		if wasteVolume > 0 {
			pww = wasteCosts / wasteVolume
		}

		pr := 0.0
		if t.volumeRainwater > 0 {
			pr = t.costsRainwater / t.volumeRainwater
		}

		fp := 0.0
		if t.fixedCount > 0 {
			fp = t.fixedPriceSum / float64(t.fixedCount)
		}

		priceWater = append(priceWater, pw)
		priceWaste = append(priceWaste, pww)
		priceRain = append(priceRain, pr)
		fixedPrices = append(fixedPrices, fp)
	}

	return models.WaterCharts{
		Consumption: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: language.T(configs.GetLanguage(), "volumeWater"),
					Data:  volWater,
					Unit:  "m³",
					YAxis: "y",
				},
				{
					Label: language.T(configs.GetLanguage(), "volumeWastewater"),
					Data:  volWaste,
					Unit:  "m³",
					YAxis: "y",
				},
				{
					Label: language.T(configs.GetLanguage(), "costsWater"),
					Data:  costWater,
					Unit:  "€",
					YAxis: "y1",
				},
				{
					Label: language.T(configs.GetLanguage(), "costsWastewater"),
					Data:  costWaste,
					Unit:  "€",
					YAxis: "y1",
				},
			},
		},

		Price: models.ChartModel{
			Labels: labels,
			Sets: []models.ChartDataset{
				{
					Label: language.T(configs.GetLanguage(), "priceWater"),
					Data:  priceWater,
					Unit:  "€/m³",
				},
				{
					Label: language.T(configs.GetLanguage(), "priceWastewater"),
					Data:  priceWaste,
					Unit:  "€/m³",
				},
				{
					Label: language.T(configs.GetLanguage(), "priceRainwater"),
					Data:  priceRain,
					Unit:  "€/m³",
				},
				{
					Label:  language.T(configs.GetLanguage(), "fixedPrice"),
					Data:   fixedPrices,
					Unit:   "€",
					Hidden: true,
				},
			},
		},
	}, nil
}
