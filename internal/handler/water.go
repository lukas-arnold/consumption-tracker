package handler

import (
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type WaterView struct {
	Water   []models.Water
	Charts  models.WaterCharts
	Summary ConsumptionSummary
}

func (h *Handler) HandleWaterView(
	w http.ResponseWriter,
	r *http.Request,
) {
	waterEntries, err :=
		h.store.GetWater()

	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err :=
		h.store.GetWaterCharts()

	if err != nil {
		handleError(w, err, 404)
		return
	}

	view := WaterView{
		Water:   waterEntries,
		Charts:  charts,
		Summary: buildWaterSummary(waterEntries),
	}

	h.renderTemplate(
		w,
		"templates/water/index.html",
		view,
	)
}

func (h *Handler) HandleAddWaterGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.renderTemplate(
		w,
		"templates/water/add.html",
		nil,
	)
}

func (h *Handler) HandleAddWaterPost(
	w http.ResponseWriter,
	r *http.Request,
) {

	year, err :=
		utils.ConvertInt(
			r.FormValue("year"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	volumeWater, err :=
		utils.ConvertFloat(
			r.FormValue("volumeWater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	volumeWastewater, err :=
		utils.ConvertFloat(
			r.FormValue("volumeWastewater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	volumeRainwater, err :=
		utils.ConvertFloat(
			r.FormValue("volumeRainwater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costsWater, err :=
		utils.ConvertFloat(
			r.FormValue("costsWater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costsWastewater, err :=
		utils.ConvertFloat(
			r.FormValue("costsWastewater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costsRainwater, err :=
		utils.ConvertFloat(
			r.FormValue("costsRainwater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	payments, err :=
		utils.ConvertFloat(
			r.FormValue("payments"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	fixedPrice, err :=
		utils.ConvertFloat(
			r.FormValue("fixedPrice"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	err =
		h.store.AddWater(
			models.WaterInput{
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
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(
		w,
		r,
		"/water",
		http.StatusFound,
	)
}

func (h *Handler) HandleEditWater(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err :=
		utils.ConvertId(
			r.PathValue("id"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry, err :=
		h.store.GetWaterEntry(id)

	if err != nil {
		handleError(w, err, 404)
		return
	}

	h.renderTemplate(
		w,
		"templates/water/edit.html",
		entry,
	)
}

func (h *Handler) HandleSaveWater(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err :=
		utils.ConvertId(
			r.PathValue("id"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	year, err :=
		utils.ConvertInt(
			r.FormValue("year"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	volumeWater, err :=
		utils.ConvertFloat(
			r.FormValue("volumeWater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	volumeWastewater, err :=
		utils.ConvertFloat(
			r.FormValue("volumeWastewater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	volumeRainwater, err :=
		utils.ConvertFloat(
			r.FormValue("volumeRainwater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costsWater, err :=
		utils.ConvertFloat(
			r.FormValue("costsWater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costsWastewater, err :=
		utils.ConvertFloat(
			r.FormValue("costsWastewater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costsRainwater, err :=
		utils.ConvertFloat(
			r.FormValue("costsRainwater"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	payments, err :=
		utils.ConvertFloat(
			r.FormValue("payments"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	fixedPrice, err :=
		utils.ConvertFloat(
			r.FormValue("fixedPrice"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	err =
		h.store.UpdateWater(
			models.Water{
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
			},
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(
		w,
		r,
		"/water",
		http.StatusFound,
	)
}

func (h *Handler) HandleDeleteWater(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err :=
		utils.ConvertId(
			r.PathValue("id"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	err =
		h.store.DeleteWater(id)

	if err != nil {
		handleError(w, err, 404)
		return
	}

	http.Redirect(
		w,
		r,
		"/water",
		http.StatusFound,
	)
}
