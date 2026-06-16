package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
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
	HandleElectricityView(w, r)
}

type ConsumptionSummary struct {
	TotalConsumption          float64
	TotalCosts                float64
	AverageConsumptionPerYear float64
	AverageCostPerUnit        float64
	YearsCount                int
}
