package handler

import (
	"html/template"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type ElectricityView struct {
	Electricity []models.Electricity
	Charts      models.ElectricityCharts
	Summary     ConsumptionSummary
}

func HandleElectricityView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				"templates/electricity/index.html",
			),
	)

	electricityEntries, err := storage.GetElectricities()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err := storage.GetElectricityCharts()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	view := ElectricityView{
		Electricity: electricityEntries,
		Charts:      charts,
		Summary:     buildElectricitySummary(electricityEntries),
	}

	if err := tmpl.Execute(w, view); err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddElectricityGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/electricity/add.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddElectricityPost(w http.ResponseWriter, r *http.Request) {
	consumption, err := utils.ConvertFloat(r.FormValue("consumption"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costs, err := utils.ConvertFloat(r.FormValue("costs"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	payments, err := utils.ConvertFloat(r.FormValue("payments"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.AddElectricity(models.ElectricityInput{
		TimeFrom:    r.FormValue("timeFrom"),
		TimeTo:      r.FormValue("timeTo"),
		Consumption: consumption,
		Costs:       costs,
		Retailer:    r.FormValue("retailer"),
		Payments:    payments,
		Note:        r.FormValue("note"),
	})
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func HandleEditElectricity(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/electricity/edit.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	entry, err := storage.GetElectricity(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}

	err = tmpl.Execute(w, entry)
	if err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleSaveElectricity(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	consumption, err := utils.ConvertFloat(r.FormValue("consumption"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costs, err := utils.ConvertFloat(r.FormValue("costs"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	payments, err := utils.ConvertFloat(r.FormValue("payments"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry := models.Electricity{
		Id: id,
		ElectricityInput: models.ElectricityInput{
			TimeFrom:    r.FormValue("timeFrom"),
			TimeTo:      r.FormValue("timeTo"),
			Consumption: consumption,
			Costs:       costs,
			Retailer:    r.FormValue("retailer"),
			Payments:    payments,
			Note:        r.FormValue("note"),
		},
	}

	err = storage.UpdateElectricity(entry)
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func HandleDeleteElectricity(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.DeleteElectricity(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}
	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func buildElectricitySummary(electricityEntries []models.Electricity) ConsumptionSummary {
	totalConsumption := 0.0
	totalCosts := 0.0

	for _, entry := range electricityEntries {
		totalConsumption += entry.Consumption
		totalCosts += entry.Costs
	}

	yearsCount := 0
	averageConsumption := 0.0

	if len(electricityEntries) > 0 {
		newest, err1 := utils.ConvertTime(electricityEntries[0].TimeTo)
		oldest, err2 := utils.ConvertTime(electricityEntries[len(electricityEntries)-1].TimeFrom)

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
