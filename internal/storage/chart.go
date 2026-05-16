package storage

import (
	"fmt"
	"sort"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
)

func GetElectricityForChart() (models.ElectricityForChart, error) {
	electricities, err := GetElectricities()
	if err != nil {
		return models.ElectricityForChart{}, err
	}
	yearTotals := map[string]struct {
		usage float64
		costs float64
	}{}
	for _, value := range electricities {
		year := value.TimeFrom
		if len(year) >= 4 {
			year = year[:4]
		}
		totals := yearTotals[year]
		totals.usage += value.Usage
		totals.costs += value.Costs
		yearTotals[year] = totals
	}
	var labels []string
	for year := range yearTotals {
		labels = append(labels, year)
	}
	sort.Strings(labels)
	var usages, costs, prices []float64
	for _, year := range labels {
		totals := yearTotals[year]
		price := 0.0
		if totals.usage > 0 {
			price = totals.costs / totals.usage
		}
		usages = append(usages, totals.usage)
		costs = append(costs, totals.costs)
		prices = append(prices, price)
	}
	return models.ElectricityForChart{Electricities: electricities, Labels: labels, Usages: usages, Costs: costs, Prices: prices}, nil
}

func GetOilForChart() (models.OilForChart, error) {
	oils, err := GetOil()
	if err != nil {
		return models.OilForChart{}, err
	}
	yearTotals := map[string]struct {
		volume float64
		costs  float64
	}{}
	for _, value := range oils {
		year := value.Date
		if len(year) >= 4 {
			year = year[:4]
		}
		totals := yearTotals[year]
		totals.volume += value.Volume
		totals.costs += value.Costs
		yearTotals[year] = totals
	}
	var labels []string
	for year := range yearTotals {
		labels = append(labels, year)
	}
	sort.Strings(labels)
	var volumes, costs, prices []float64
	for _, year := range labels {
		totals := yearTotals[year]
		price := 0.0
		if totals.volume > 0 {
			price = totals.costs / totals.volume
		}
		volumes = append(volumes, totals.volume)
		costs = append(costs, totals.costs)
		prices = append(prices, price)
	}
	return models.OilForChart{Oils: oils, Labels: labels, Volumes: volumes, Costs: costs, Prices: prices}, nil
}

func GetOilFillLevelForChart() (models.OilFillLevelForChart, error) {
	fillLevels, err := GetOilFillLevels()
	if err != nil {
		return models.OilFillLevelForChart{}, err
	}
	var dates []string
	var levels []float64
	for _, value := range fillLevels {
		dates = append(dates, value.Date)
		levels = append(levels, value.Level)
	}
	return models.OilFillLevelForChart{OilFillLevels: fillLevels, Dates: dates, Levels: levels}, nil
}

func GetWaterForChart() (models.WaterForChart, error) {
	waters, err := GetWater()
	if err != nil {
		return models.WaterForChart{}, err
	}
	// Aggregate per-year totals separately for water, wastewater and rainwater
	yearTotals := map[int]struct {
		volumeWater     float64
		volumeWastewater float64
		volumeRainwater float64
		costsWater      float64
		costsWastewater float64
		costsRainwater  float64
		totalFixedPrice float64
		fixedCount      int
	}{}
	for _, v := range waters {
		t := yearTotals[v.Year]
		t.volumeWater += v.VolumeWater
		t.volumeWastewater += v.VolumeWastewater
		t.volumeRainwater += v.VolumeRainwater
		t.costsWater += v.CostsWater
		t.costsWastewater += v.CostsWastewater
		t.costsRainwater += v.CostsRainwater
		if v.FixedPrice > 0 {
			t.totalFixedPrice += v.FixedPrice
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
	var volumes, volumesWater, volumesWastewater, volumesRainwater []float64
	var costs, costsWater, costsWastewater, costsRainwater []float64
	var prices []float64
	var pricesWater, pricesWastewater, pricesRainwater, fixedPrices []float64
	for _, y := range years {
		labels = append(labels, fmt.Sprintf("%d", y))
		t := yearTotals[y]
		totalVolume := t.volumeWater + t.volumeWastewater + t.volumeRainwater
		totalCosts := t.costsWater + t.costsWastewater + t.costsRainwater

		// compute per-type prices (cost / volume) when possible
		priceWater := 0.0
		if t.volumeWater > 0 {
			priceWater = t.costsWater / t.volumeWater
		}
		priceWastewater := 0.0
		if t.volumeWastewater > 0 {
			priceWastewater = t.costsWastewater / t.volumeWastewater
		}
		priceRainwater := 0.0
		if t.volumeRainwater > 0 {
			priceRainwater = t.costsRainwater / t.volumeRainwater
		}
		fixedPrice := 0.0
		if t.fixedCount > 0 {
			fixedPrice = t.totalFixedPrice / float64(t.fixedCount)
		}

		volumes = append(volumes, totalVolume)
		volumesWater = append(volumesWater, t.volumeWater)
		volumesWastewater = append(volumesWastewater, t.volumeWastewater)
		volumesRainwater = append(volumesRainwater, t.volumeRainwater)

		costs = append(costs, totalCosts)
		costsWater = append(costsWater, t.costsWater)
		costsWastewater = append(costsWastewater, t.costsWastewater)
		costsRainwater = append(costsRainwater, t.costsRainwater)

		prices = append(prices, 0) // legacy/unused combined price kept as placeholder
		pricesWater = append(pricesWater, priceWater)
		pricesWastewater = append(pricesWastewater, priceWastewater)
		pricesRainwater = append(pricesRainwater, priceRainwater)
		fixedPrices = append(fixedPrices, fixedPrice)
	}

	return models.WaterForChart{
		Waters: waters,
		Labels: labels,
		Volumes: volumes,
		VolumesWater: volumesWater,
		VolumesWastewater: volumesWastewater,
		VolumesRainwater: volumesRainwater,
		Costs: costs,
		CostsWater: costsWater,
		CostsWastewater: costsWastewater,
		CostsRainwater: costsRainwater,
		Prices: prices,
		PricesWater: pricesWater,
		PricesWastewater: pricesWastewater,
		PricesRainwater: pricesRainwater,
		FixedPrices: fixedPrices,
	}, nil
}
