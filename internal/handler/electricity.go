package handler

import (
	"html/template"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type ElectricityView struct {
	Electricity []models.Electricity
	Charts      models.ElectricityCharts
	Summary     ConsumptionSummary
}

func (h *Handler) HandleElectricityView(
	w http.ResponseWriter,
	r *http.Request,
) {

	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				"templates/electricity/index.html",
			),
	)

	electricityEntries, err :=
		h.store.GetElectricities()

	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err :=
		h.store.GetElectricityCharts()

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

func (h *Handler) HandleAddElectricityGet(
	w http.ResponseWriter,
	r *http.Request,
) {

	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				"templates/electricity/add.html",
			),
	)

	if err := tmpl.Execute(w, nil); err != nil {
		handleError(w, err, 500)
		return
	}
}

func (h *Handler) HandleAddElectricityPost(
	w http.ResponseWriter,
	r *http.Request,
) {

	consumption, err :=
		utils.ConvertFloat(
			r.FormValue("consumption"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costs, err :=
		utils.ConvertFloat(
			r.FormValue("costs"),
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

	err = h.store.AddElectricity(
		models.ElectricityInput{
			TimeFrom:    r.FormValue("timeFrom"),
			TimeTo:      r.FormValue("timeTo"),
			Consumption: consumption,
			Costs:       costs,
			Retailer:    r.FormValue("retailer"),
			Payments:    payments,
			Note:        r.FormValue("note"),
		},
	)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(
		w,
		r,
		"/electricity",
		http.StatusFound,
	)
}

func (h *Handler) HandleEditElectricity(
	w http.ResponseWriter,
	r *http.Request,
) {

	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				"templates/electricity/edit.html",
			),
	)

	id, err :=
		utils.ConvertId(
			r.PathValue("id"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry, err :=
		h.store.GetElectricity(id)

	if err != nil {
		handleError(w, err, 404)
		return
	}

	if err := tmpl.Execute(w, entry); err != nil {
		handleError(w, err, 500)
		return
	}
}

func (h *Handler) HandleSaveElectricity(
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

	consumption, err :=
		utils.ConvertFloat(
			r.FormValue("consumption"),
		)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	costs, err :=
		utils.ConvertFloat(
			r.FormValue("costs"),
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

	err = h.store.UpdateElectricity(entry)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(
		w,
		r,
		"/electricity",
		http.StatusFound,
	)
}

func (h *Handler) HandleDeleteElectricity(
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

	err = h.store.DeleteElectricity(id)

	if err != nil {
		handleError(w, err, 404)
		return
	}

	http.Redirect(
		w,
		r,
		"/electricity",
		http.StatusFound,
	)
}
