package storage

import (
	"fmt"
	"sort"
	"time"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func (s *Storage) GetElectricityCharts() (models.ElectricityCharts, error) {
	electricities, err := s.GetElectricities()
	if err != nil {
		return models.ElectricityCharts{}, err
	}

	type totals struct{ consumption, costs float64 }
	yearTotals := map[string]totals{}

	for _, v := range electricities {
		from, err1 := time.Parse("2006-01-02", v.TimeFrom)
		to, err2 := time.Parse("2006-01-02", v.TimeTo)
		if err1 != nil || err2 != nil {
			return models.ElectricityCharts{}, err1
		}

		totalDays := to.Sub(from).Hours() / 24
		if totalDays <= 0 {
			continue
		}

		for year := from.Year(); year <= to.Year(); year++ {
			yStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
			yEnd := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)

			pStart := from
			if yStart.After(pStart) {
				pStart = yStart
			}

			pEnd := to
			if yEnd.Before(pEnd) {
				pEnd = yEnd
			}

			daysInYear := pEnd.Sub(pStart).Hours() / 24
			if daysInYear <= 0 {
				continue
			}

			key := fmt.Sprintf("%d", year)
			t := yearTotals[key]
			ratio := daysInYear / totalDays
			t.consumption += v.Consumption * ratio
			t.costs += v.Costs * ratio
			yearTotals[key] = t
		}
	}

	labels := sortedKeys(yearTotals)
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

	lang := configs.GetLanguage()
	return models.ElectricityCharts{
		Consumption: models.ChartModel{Labels: labels, Sets: []models.ChartDataset{
			{Label: language.T(lang, "consumption"), Data: consumption, Unit: "kWh", YAxis: "y"},
			{Label: language.T(lang, "costs"), Data: costs, Unit: "€", YAxis: "y1"},
		}},
		Price: models.ChartModel{Labels: labels, Sets: []models.ChartDataset{
			{Label: language.T(lang, "pricePerUnit"), Data: prices, Unit: "€/kWh"},
		}},
	}, nil
}

// GetOilCharts aggregates oil consumption and fill levels.
func (s *Storage) GetOilCharts() (models.OilCharts, error) {
	oils, err := s.GetOil()
	if err != nil {
		return models.OilCharts{}, err
	}

	type totals struct{ volume, costs float64 }
	yearTotals := map[string]totals{}

	for _, v := range oils {
		year := v.Date[:4]
		t := yearTotals[year]
		t.volume += v.Volume
		t.costs += v.Costs
		yearTotals[year] = t
	}

	labels := sortedKeys(yearTotals)
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

	fillLevels, _ := s.GetOilFillLevels()
	sort.Slice(fillLevels, func(i, j int) bool { return fillLevels[i].Date < fillLevels[j].Date })

	var fillDates []string
	var fillValues, fillPercentages []float64
	for _, v := range fillLevels {
		fillDates = append(fillDates, v.Date)
		fillValues = append(fillValues, v.Level)
		fillPercentages = append(fillPercentages, v.Percentage)
	}

	lang := configs.GetLanguage()
	return models.OilCharts{
		Consumption: models.ChartModel{Labels: labels, Sets: []models.ChartDataset{
			{Label: language.T(lang, "volume"), Data: volumes, Unit: "l", YAxis: "y"},
			{Label: language.T(lang, "costs"), Data: costs, Unit: "€", YAxis: "y1"},
		}},
		Price: models.ChartModel{Labels: labels, Sets: []models.ChartDataset{
			{Label: language.T(lang, "pricePerUnit"), Data: prices, Unit: "€/l"},
		}},
		FillLevel: models.ChartModel{Labels: fillDates, Sets: []models.ChartDataset{
			{Label: language.T(lang, "levelCm"), Data: fillValues, Unit: "cm", YAxis: "y"},
			{Label: language.T(lang, "level%"), Data: fillPercentages, Unit: "%", YAxis: "y1"},
		}},
	}, nil
}

func (s *Storage) GetWaterCharts() (models.WaterCharts, error) {
	waters, err := s.GetWater()
	if err != nil {
		return models.WaterCharts{}, err
	}

	type totals struct {
		volW, volWW, volR    float64
		costW, costWW, costR float64
		fixedPriceSum        float64
		fixedCount           int
	}

	yearTotals := map[int]totals{}
	for _, v := range waters {
		t := yearTotals[v.Year]
		t.volW += v.VolumeWater
		t.volWW += v.VolumeWastewater
		t.volR += v.VolumeRainwater
		t.costW += v.CostsWater
		t.costWW += v.CostsWastewater
		t.costR += v.CostsRainwater
		if v.FixedPrice > 0 {
			t.fixedPriceSum += v.FixedPrice
			t.fixedCount++
		}
		yearTotals[v.Year] = t
	}

	years := make([]int, 0, len(yearTotals))
	for y := range yearTotals {
		years = append(years, y)
	}
	sort.Ints(years)

	var labels []string
	var volW, volWW, costW, costWW []float64
	var priceW, priceWW, priceR, fixedPrices []float64

	for _, y := range years {
		t := yearTotals[y]
		labels = append(labels, fmt.Sprintf("%d", y))

		volW = append(volW, t.volW)
		volWW = append(volWW, t.volWW+t.volR)
		costW = append(costW, t.costW)
		costWW = append(costWW, t.costWW+t.costR)

		priceW = append(priceW, safeDiv(t.costW, t.volW))
		priceWW = append(priceWW, safeDiv(t.costWW, t.volWW))
		priceR = append(priceR, safeDiv(t.costR, t.volR))

		fp := 0.0
		if t.fixedCount > 0 {
			fp = t.fixedPriceSum / float64(t.fixedCount)
		}
		fixedPrices = append(fixedPrices, fp)
	}

	lang := configs.GetLanguage()
	return models.WaterCharts{
		Consumption: models.ChartModel{Labels: labels, Sets: []models.ChartDataset{
			{Label: language.T(lang, "volumeWater"), Data: volW, Unit: "m³", YAxis: "y"},
			{Label: language.T(lang, "volumeWastewater"), Data: volWW, Unit: "m³", YAxis: "y"},
			{Label: language.T(lang, "costsWater"), Data: costW, Unit: "€", YAxis: "y1"},
			{Label: language.T(lang, "costsWastewater"), Data: costWW, Unit: "€", YAxis: "y1"},
		}},
		Price: models.ChartModel{Labels: labels, Sets: []models.ChartDataset{
			{Label: language.T(lang, "priceWater"), Data: priceW, Unit: "€/m³"},
			{Label: language.T(lang, "priceWastewater"), Data: priceWW, Unit: "€/m³"},
			{Label: language.T(lang, "priceRainwater"), Data: priceR, Unit: "€/m³"},
			{Label: language.T(lang, "fixedPrice"), Data: fixedPrices, Unit: "€", Hidden: true},
		}},
	}, nil
}

func safeDiv(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
