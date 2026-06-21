package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func handleError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("HTTP %d: %v", statusCode, err)
	http.Error(w, http.StatusText(statusCode), statusCode)
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

func calculateMonthlyPeriod(payments, days float64) float64 {
	return payments / days * 30.4375 // average days per month
}

func formatMonthly(payments float64, timeFrom, timeTo string) string {
	if payments == 0 {
		return "-"
	}

	// Assume the payment amount covers a full year when:
	// - no period information is available
	// - the dates cannot be parsed
	// - the calculated period is invalid
	// - the period is approximately one calendar year
	monthly := payments / 12

	if timeFrom == "" || timeTo == "" {
		return formatEuro(monthly)
	}

	from, err1 := utils.ConvertTime(timeFrom)
	to, err2 := utils.ConvertTime(timeTo)

	if err1 != nil || err2 != nil {
		return formatEuro(monthly)
	}

	days := to.Sub(from).Hours() / 24

	if days <= 0 || (days >= 364 && days <= 366) {
		return formatEuro(monthly)
	}

	return formatEuro(calculateMonthlyPeriod(payments, days))
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

func HandleFiles(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/web/", http.FileServerFS(configs.GetWebFiles())).ServeHTTP(w, r)
}

func HandleServiceWorker(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, configs.GetWebFiles(), "service-worker.js")
}

func HandleView(w http.ResponseWriter, r *http.Request) {
	HandleElectricityView(w, r)
}

type ConsumptionSummary struct {
	TotalConsumption          float64
	TotalCosts                float64
	AverageConsumptionPerYear float64
	AverageCostPerUnit        float64
	YearsCount                int
}
