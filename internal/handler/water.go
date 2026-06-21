package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type WaterView struct {
	Water   []models.Water
	Charts  models.WaterCharts
	Summary ConsumptionSummary
}

func HandleWaterView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				"templates/water/index.html",
			),
	)

	waterEntries, err := storage.GetWater()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err := storage.GetWaterCharts()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	view := WaterView{
		Water:   waterEntries,
		Charts:  charts,
		Summary: buildWaterSummary(waterEntries),
	}

	if err := tmpl.Execute(w, view); err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddWaterGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/water/add.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		handleError(w, err, 500)
		return
	}
}

func HandleAddWaterPost(w http.ResponseWriter, r *http.Request) {
	year, err := utils.ConvertInt(r.FormValue("year"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volumeWater, err := utils.ConvertFloat(r.FormValue("volumeWater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volumeWastewater, err := utils.ConvertFloat(r.FormValue("volumeWastewater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volumeRainwater, err := utils.ConvertFloat(r.FormValue("volumeRainwater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costsWater, err := utils.ConvertFloat(r.FormValue("costsWater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costsWastewater, err := utils.ConvertFloat(r.FormValue("costsWastewater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costsRainwater, err := utils.ConvertFloat(r.FormValue("costsRainwater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	payments, err := utils.ConvertFloat(r.FormValue("payments"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	fixedPrice, err := utils.ConvertFloat(r.FormValue("fixedPrice"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.AddWater(models.WaterInput{
		Year:             year,
		VolumeWater:      volumeWater,
		VolumeWastewater: volumeWastewater,
		VolumeRainwater:  volumeRainwater,
		CostsWater:       costsWater,
		CostsWastewater:  costsWastewater,
		CostsRainwater:   costsRainwater,
		Payments:         payments,
		FixedPrice:       fixedPrice,
		Note:             r.FormValue("note"),
	})
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/water", http.StatusFound)
}

func HandleEditWater(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/water/edit.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	entry, err := storage.GetWaterEntry(id)
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

func HandleSaveWater(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	year, err := utils.ConvertInt(r.FormValue("year"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volumeWater, err := utils.ConvertFloat(r.FormValue("volumeWater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volumeWastewater, err := utils.ConvertFloat(r.FormValue("volumeWastewater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	volumeRainwater, err := utils.ConvertFloat(r.FormValue("volumeRainwater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costsWater, err := utils.ConvertFloat(r.FormValue("costsWater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costsWastewater, err := utils.ConvertFloat(r.FormValue("costsWastewater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	costsRainwater, err := utils.ConvertFloat(r.FormValue("costsRainwater"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	payments, err := utils.ConvertFloat(r.FormValue("payments"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	fixedPrice, err := utils.ConvertFloat(r.FormValue("fixedPrice"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	entry := models.Water{
		Id: id,
		WaterInput: models.WaterInput{
			Year:             year,
			VolumeWater:      volumeWater,
			VolumeWastewater: volumeWastewater,
			VolumeRainwater:  volumeRainwater,
			CostsWater:       costsWater,
			CostsWastewater:  costsWastewater,
			CostsRainwater:   costsRainwater,
			Payments:         payments,
			FixedPrice:       fixedPrice,
			Note:             r.FormValue("note"),
		},
	}

	err = storage.UpdateWater(entry)
	if err != nil {
		handleError(w, err, 500)
		return
	}
	http.Redirect(w, r, "/water", http.StatusFound)
}

func HandleDeleteWater(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}
	err = storage.DeleteWater(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}
	http.Redirect(w, r, "/water", http.StatusFound)
}

func buildWaterSummary(waterEntries []models.Water) ConsumptionSummary {
	years := map[string]bool{}
	totalVolume := 0.0
	totalCosts := 0.0
	for _, entry := range waterEntries {
		years[fmt.Sprint(entry.Year)] = true
		totalVolume += entry.VolumeWater
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
