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

// Helper to reduce repetition in template parsing
func parseElectricityTemplate(file string) *template.Template {
	return template.Must(template.New("base.html").
		Funcs(getTemplateFuncs()).
		ParseFS(configs.GetWebFiles(), "templates/base.html", file))
}

func (h *Handler) HandleElectricityView(w http.ResponseWriter, r *http.Request) {
	tmpl := parseElectricityTemplate("templates/electricity/index.html")

	electricityEntries, err := h.store.GetElectricities()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err := h.store.GetElectricityCharts()
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
	}
}

func (h *Handler) HandleAddElectricityGet(w http.ResponseWriter, r *http.Request) {
	tmpl := parseElectricityTemplate("templates/electricity/add.html")

	if err := tmpl.Execute(w, nil); err != nil {
		handleError(w, err, 500)
	}
}

func (h *Handler) HandleAddElectricityPost(w http.ResponseWriter, r *http.Request) {
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

	input := models.ElectricityInput{
		TimeFrom:    r.FormValue("timeFrom"),
		TimeTo:      r.FormValue("timeTo"),
		Consumption: consumption,
		Costs:       costs,
		Retailer:    r.FormValue("retailer"),
		Payments:    payments,
		Note:        r.FormValue("note"),
	}

	if err := h.store.AddElectricity(input); err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func (h *Handler) HandleEditElectricity(w http.ResponseWriter, r *http.Request) {
	tmpl := parseElectricityTemplate("templates/electricity/edit.html")

	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry, err := h.store.GetElectricity(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}

	if err := tmpl.Execute(w, entry); err != nil {
		handleError(w, err, 500)
	}
}

func (h *Handler) HandleSaveElectricity(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.UpdateElectricity(entry); err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func (h *Handler) HandleDeleteElectricity(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	if err := h.store.DeleteElectricity(id); err != nil {
		handleError(w, err, 404)
		return
	}

	http.Redirect(w, r, "/electricity", http.StatusFound)
}
