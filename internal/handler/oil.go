package handler

import (
	"html/template"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type OilView struct {
	Oil          []models.Oil
	OilFillLevel []models.OilFillLevel
	Charts       models.OilCharts
	Summary      ConsumptionSummary
}

func HandleOilView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				"templates/oil/index.html",
			),
	)

	oilEntries, err := storage.GetOil()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	oilFillLevels, err := storage.GetOilFillLevels()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err := storage.GetOilCharts()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	view := OilView{
		Oil:          oilEntries,
		OilFillLevel: oilFillLevels,
		Charts:       charts,
		Summary:      buildOilSummary(oilEntries),
	}

	if err := tmpl.Execute(w, view); err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddOilGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/oil/add.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddOilPost(w http.ResponseWriter, r *http.Request) {
	volume, err := utils.ConvertFloat(r.FormValue("volume"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costs, err := utils.ConvertFloat(r.FormValue("costs"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.AddOil(models.OilInput{
		Date:     r.FormValue("date"),
		Volume:   volume,
		Costs:    costs,
		Retailer: r.FormValue("retailer"),
		Note:     r.FormValue("note"),
	})
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleEditOil(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/oil/edit.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	entry, err := storage.GetOilEntry(id)
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

func HandleSaveOil(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volume, err := utils.ConvertFloat(r.FormValue("volume"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costs, err := utils.ConvertFloat(r.FormValue("costs"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry := models.Oil{
		Id: id,
		OilInput: models.OilInput{
			Date:     r.FormValue("date"),
			Volume:   volume,
			Costs:    costs,
			Retailer: r.FormValue("retailer"),
			Note:     r.FormValue("note"),
		},
	}

	err = storage.UpdateOil(entry)
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleDeleteOil(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.DeleteOil(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleAddOilFillLevelGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/oil/addFillLevel.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddOilFillLevelPost(w http.ResponseWriter, r *http.Request) {
	level, err := utils.ConvertFloat(r.FormValue("level"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.AddOilFillLevel(models.OilFillLevelInput{
		Date:  r.FormValue("date"),
		Level: level,
	})
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleEditOilFillLevel(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/oil/editFillLevel.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	entry, err := storage.GetOilFillLevel(id)
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

func HandleSaveOilFillLevel(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	level, err := utils.ConvertFloat(r.FormValue("level"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry := models.OilFillLevel{
		Id: id,
		OilFillLevelInput: models.OilFillLevelInput{
			Date:  r.FormValue("date"),
			Level: level,
		},
	}

	err = storage.UpdateOilFillLevel(entry)
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleDeleteOilFillLevel(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.DeleteOilFillLevel(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func buildOilSummary(oilEntries []models.Oil) ConsumptionSummary {
	totalVolume := 0.0
	totalCosts := 0.0

	for _, entry := range oilEntries {
		totalVolume += entry.Volume
		totalCosts += entry.Costs
	}

	yearsCount := 0
	averageVolume := 0.0

	if len(oilEntries) > 0 {
		newest, err1 := utils.ConvertTime(oilEntries[0].Date)
		oldest, err2 := utils.ConvertTime(oilEntries[len(oilEntries)-1].Date)

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
