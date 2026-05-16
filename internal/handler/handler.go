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
	return formatEuro(payments/12)
}

func formatDifference(payments, costs float64) string {
	difference := payments - costs
	if difference == 0 {
		return "-"
	}
	return formatEuro(difference)
}

func formatTotalCosts(costsWater, costsWastewater, costsRainwater float64) string {
	total := costsWater + costsWastewater + costsRainwater
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
		"Float": formatFloat,
		"Int":   formatInt,
		"Euro":  formatEuro,
		"Price": formatPrice,
		"Monthly": formatMonthly,
		"Difference": formatDifference,
		"TotalCosts": formatTotalCosts,
		"Add": add,
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

func handleElectricityView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("view.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/electricity/view.html", "templates/electricity/index.html"),
	)
	usageStorage, err := storage.GetUsageStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	chart, err := storage.GetElectricityForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	summary := buildElectricitySummary(usageStorage)
	view := struct {
		models.UsageStorage
		Chart   models.ElectricityForChart
		Summary ConsumptionSummary
	}{usageStorage, chart, summary}
	err = tmpl.Execute(w, view)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleOilView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("view.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/oil/view.html", "templates/oil/index.html"),
	)
	usageStorage, err := storage.GetUsageStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	fillLevels, err := storage.GetOilFillLevels()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	usageStorage.OilFillLevels = fillLevels
	chart, err := storage.GetOilForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	fillChart, err := storage.GetOilFillLevelForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	summary := buildOilSummary(usageStorage)
	view := struct {
		models.UsageStorage
		OilChart       models.OilForChart
		FillLevelChart models.OilFillLevelForChart
		Summary        ConsumptionSummary
	}{usageStorage, chart, fillChart, summary}
	err = tmpl.Execute(w, view)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

type ConsumptionSummary struct {
	TotalUsage          float64
	TotalCosts          float64
	AverageUsagePerYear float64
	AverageCostPerUnit  float64
	YearsCount          int
}

func HandleWaterView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("view.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/water/view.html", "templates/water/index.html"),
	)
	usageStorage, err := storage.GetUsageStorage()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	chart, err := storage.GetWaterForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	summary := buildWaterSummary(usageStorage)
	view := struct {
		models.UsageStorage
		Chart   models.WaterForChart
		Summary ConsumptionSummary
	}{usageStorage, chart, summary}
	err = tmpl.Execute(w, view)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func buildElectricitySummary(usageStorage models.UsageStorage) ConsumptionSummary {
	years := map[string]bool{}
	totalUsage := 0.0
	totalCosts := 0.0
	for _, entry := range usageStorage.Electricity {
		if len(entry.TimeFrom) >= 4 {
			years[entry.TimeFrom[:4]] = true
		}
		totalUsage += entry.Usage
		totalCosts += entry.Costs
	}
	yearsCount := len(years)
	averageUsage := 0.0
	if yearsCount > 0 {
		averageUsage = totalUsage / float64(yearsCount)
	}
	averageCost := 0.0
	if totalUsage > 0 {
		averageCost = totalCosts / totalUsage
	}
	return ConsumptionSummary{TotalUsage: totalUsage, TotalCosts: totalCosts, AverageUsagePerYear: averageUsage, AverageCostPerUnit: averageCost, YearsCount: yearsCount}
}

func buildOilSummary(usageStorage models.UsageStorage) ConsumptionSummary {
	years := map[string]bool{}
	totalVolume := 0.0
	totalCosts := 0.0
	for _, entry := range usageStorage.Oil {
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
	return ConsumptionSummary{TotalUsage: totalVolume, TotalCosts: totalCosts, AverageUsagePerYear: averageVolume, AverageCostPerUnit: averageCost, YearsCount: yearsCount}
}

func buildWaterSummary(usageStorage models.UsageStorage) ConsumptionSummary {
	years := map[string]bool{}
	totalVolume := 0.0
	totalCosts := 0.0
	for _, entry := range usageStorage.Water {
		years[fmt.Sprint(entry.Year)] = true
		totalVolume += entry.VolumeWater + entry.VolumeWastewater + entry.VolumeRainwater
		totalCosts += entry.CostsWater + entry.CostsWastewater + entry.CostsRainwater
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
		TotalUsage:          totalVolume,
		TotalCosts:          totalCosts,
		AverageUsagePerYear: averageVolume,
		AverageCostPerUnit:  averageCost,
		YearsCount:          yearsCount,
	}
}

func HandleElectricityHistory(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("chart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/electricity/chart.html"),
	)
	chart, err := storage.GetElectricityForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleElectricityChart(w http.ResponseWriter, r *http.Request) {
	HandleElectricityHistory(w, r)
}

func HandleOilHistory(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("chart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/oil/chart.html"),
	)
	chart, err := storage.GetOilForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleOilChart(w http.ResponseWriter, r *http.Request) {
	HandleOilHistory(w, r)
}

func HandleOilFillLevelHistory(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("chart.html").Funcs(getTemplateFuncs()).ParseFS(configs.GetWebFiles(), "templates/oil/fillchart.html"),
	)
	chart, err := storage.GetOilFillLevelForChart()
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
	}
	err = tmpl.Execute(w, chart)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleOilFillLevelChart(w http.ResponseWriter, r *http.Request) {
	HandleOilFillLevelHistory(w, r)
}

// Water history/chart handlers removed — charts are embedded in the water index page.
