package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
)

func errorHandling(w http.ResponseWriter, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
}

func formatFloat(value float64, digits int) string {
	formatted := fmt.Sprintf("%.*f", digits, value)
	return strings.ReplaceAll(formatted, ".", ",")
}

func formatInt(value float64) string {
	formatted := fmt.Sprintf("%.0f", value)
	return strings.ReplaceAll(formatted, ".", ",")
}

func formatEuro(value float64) string {
	if value == 0 {
		return "-"
	}
	return formatFloat(value, 2) + " €"
}

func formatPrice(costs, units float64, unit string) string {
	if units <= 0 {
		return "-"
	}
	return formatFloat(costs/units, 3) + " €/" + unit
}

func formatMonthly(payments float64) string {
	if payments == 0 {
		return "-"
	}
	return formatEuro(payments / 12)
}

func formatDifference(payments, costs float64) string {
	difference := payments - costs
	if difference == 0 {
		return "-"
	}
	return formatEuro(difference)
}

func formatTotalCosts(costsWater, costsWastewater, costsRainwater, FixedPrice float64) string {
	total := costsWater + costsWastewater + costsRainwater + FixedPrice
	if total == 0 {
		return "-"
	}
	return formatEuro(total)
}

func add(values ...float64) float64 {
	total := 0.0
	for _, value := range values {
		total += value
	}
	return total
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"T": func(key string) string {
			return language.T(configs.GetLanguage(), key)
		},
		"Float":      formatFloat,
		"Int":        formatInt,
		"Euro":       formatEuro,
		"Price":      formatPrice,
		"Monthly":    formatMonthly,
		"Difference": formatDifference,
		"TotalCosts": formatTotalCosts,
		"Add":        add,
	}
}

func HandleServiceWorker(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, configs.GetWebFiles(), "service-worker.js")
}

func HandleFiles(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/web/", http.FileServerFS(configs.GetWebFiles())).ServeHTTP(w, r)
}

func HandleView(w http.ResponseWriter, r *http.Request) {
	handleElectricityView(w, r)
}

func HandleElectricityView(w http.ResponseWriter, r *http.Request) {
	handleElectricityView(w, r)
}

type ElectricityView struct {
	models.ConsumptionStorage
	Charts  models.ElectricityCharts
	Summary ConsumptionSummary
}

func handleElectricityView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("view.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/electricity/view.html",
				"templates/electricity/index.html",
			),
	)

	consumptionStorage, err := storage.GetConsumptionStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	charts, err := storage.GetElectricityCharts()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	view := ElectricityView{
		ConsumptionStorage: consumptionStorage,
		Charts:             charts,
		Summary:            buildElectricitySummary(consumptionStorage),
	}

	if err := tmpl.Execute(w, view); err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

type OilView struct {
	models.ConsumptionStorage
	Charts  models.OilCharts
	Summary ConsumptionSummary
}

func HandleOilView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("view.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/oil/view.html",
				"templates/oil/index.html",
			),
	)

	consumptionStorage, err := storage.GetConsumptionStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	fillLevels, err := storage.GetOilFillLevels()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	consumptionStorage.OilFillLevels = fillLevels

	charts, err := storage.GetOilCharts()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	view := OilView{
		ConsumptionStorage: consumptionStorage,
		Charts:             charts,
		Summary:            buildOilSummary(consumptionStorage),
	}

	if err := tmpl.Execute(w, view); err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

type ConsumptionSummary struct {
	TotalConsumption          float64
	TotalCosts                float64
	AverageConsumptionPerYear float64
	AverageCostPerUnit        float64
	YearsCount                int
}

type WaterView struct {
	models.ConsumptionStorage
	Charts  models.WaterCharts
	Summary ConsumptionSummary
}

func HandleWaterView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("view.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/water/view.html",
				"templates/water/index.html",
			),
	)

	consumptionStorage, err := storage.GetConsumptionStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	charts, err := storage.GetWaterCharts()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	view := WaterView{
		ConsumptionStorage: consumptionStorage,
		Charts:             charts,
		Summary:            buildWaterSummary(consumptionStorage),
	}

	if err := tmpl.Execute(w, view); err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func buildElectricitySummary(consumptionStorage models.ConsumptionStorage) ConsumptionSummary {
	years := map[string]bool{}
	totalConsumption := 0.0
	totalCosts := 0.0
	for _, entry := range consumptionStorage.Electricity {
		if len(entry.TimeFrom) >= 4 {
			years[entry.TimeFrom[:4]] = true
		}
		totalConsumption += entry.Consumption
		totalCosts += entry.Costs
	}
	yearsCount := len(years)
	averageConsumption := 0.0
	if yearsCount > 0 {
		averageConsumption = totalConsumption / float64(yearsCount)
	}
	averageCost := 0.0
	if totalConsumption > 0 {
		averageCost = totalCosts / totalConsumption
	}
	return ConsumptionSummary{TotalConsumption: totalConsumption, TotalCosts: totalCosts, AverageConsumptionPerYear: averageConsumption, AverageCostPerUnit: averageCost, YearsCount: yearsCount}
}

func buildOilSummary(consumptionStorage models.ConsumptionStorage) ConsumptionSummary {
	years := map[string]bool{}
	totalVolume := 0.0
	totalCosts := 0.0
	for _, entry := range consumptionStorage.Oil {
		if len(entry.Date) >= 4 {
			years[entry.Date[:4]] = true
		}
		totalVolume += entry.Volume
		totalCosts += entry.Costs
	}
	yearsCount := len(years)
	averageVolume := 0.0
	if yearsCount > 0 {
		averageVolume = totalVolume / float64(yearsCount)
	}
	averageCost := 0.0
	if totalVolume > 0 {
		averageCost = totalCosts / totalVolume
	}
	return ConsumptionSummary{TotalConsumption: totalVolume, TotalCosts: totalCosts, AverageConsumptionPerYear: averageVolume, AverageCostPerUnit: averageCost, YearsCount: yearsCount}
}

func buildWaterSummary(consumptionStorage models.ConsumptionStorage) ConsumptionSummary {
	years := map[string]bool{}
	totalVolume := 0.0
	totalCosts := 0.0
	for _, entry := range consumptionStorage.Water {
		years[fmt.Sprint(entry.Year)] = true
		totalVolume += entry.VolumeWater + entry.VolumeWastewater + entry.VolumeRainwater
		totalCosts += entry.CostsWater + entry.CostsWastewater + entry.CostsRainwater + entry.FixedPrice
	}
	yearsCount := len(years)
	averageVolume := 0.0
	if yearsCount > 0 {
		averageVolume = totalVolume / float64(yearsCount)
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
