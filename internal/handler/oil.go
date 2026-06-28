package handler

import (
	"html/template"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type OilView struct {
	Oil          []models.Oil
	OilFillLevel []models.OilFillLevel
	Charts       models.OilCharts
	Summary      ConsumptionSummary
}

func (h *Handler) HandleOilView(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("base.html").
		Funcs(getTemplateFuncs()).
		ParseFS(configs.GetWebFiles(), "templates/base.html", "templates/oil/index.html"))

	oilEntries, err := h.store.GetOil()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	fillLevels, err := h.store.GetOilFillLevels()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	charts, err := h.store.GetOilCharts()
	if err != nil {
		handleError(w, err, 404)
		return
	}

	view := OilView{
		Oil:          oilEntries,
		OilFillLevel: fillLevels,
		Charts:       charts,
		Summary:      buildOilSummary(oilEntries),
	}

	if err := tmpl.Execute(w, view); err != nil {
		handleError(w, err, 500)
	}
}

func (h *Handler) HandleAddOilGet(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "templates/oil/add.html", nil)
}

func (h *Handler) HandleAddOilPost(w http.ResponseWriter, r *http.Request) {
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

	input := models.OilInput{
		Date:     r.FormValue("date"),
		Volume:   volume,
		Costs:    costs,
		Retailer: r.FormValue("retailer"),
		Note:     r.FormValue("note"),
	}

	if err := h.store.AddOil(input); err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(w, r, "/oil", http.StatusFound)
}

func (h *Handler) HandleEditOil(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry, err := h.store.GetOilEntry(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}

	h.renderTemplate(w, "templates/oil/edit.html", entry)
}

func (h *Handler) HandleSaveOil(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.UpdateOil(entry); err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(w, r, "/oil", http.StatusFound)
}

func (h *Handler) HandleDeleteOil(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	if err := h.store.DeleteOil(id); err != nil {
		handleError(w, err, 404)
		return
	}

	http.Redirect(w, r, "/oil", http.StatusFound)
}

func (h *Handler) HandleAddOilFillLevelGet(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "templates/oil/addFillLevel.html", nil)
}

func (h *Handler) HandleAddOilFillLevelPost(w http.ResponseWriter, r *http.Request) {
	level, err := utils.ConvertFloat(r.FormValue("level"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	input := models.OilFillLevelInput{
		Date:  r.FormValue("date"),
		Level: level,
	}

	if err := h.store.AddOilFillLevel(input); err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(w, r, "/oil", http.StatusFound)
}

func (h *Handler) HandleEditOilFillLevel(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	entry, err := h.store.GetOilFillLevel(id)
	if err != nil {
		handleError(w, err, 404)
		return
	}

	h.renderTemplate(w, "templates/oil/editFillLevel.html", entry)
}

func (h *Handler) HandleSaveOilFillLevel(w http.ResponseWriter, r *http.Request) {
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

	if err := h.store.UpdateOilFillLevel(entry); err != nil {
		handleError(w, err, 500)
		return
	}

	http.Redirect(w, r, "/oil", http.StatusFound)
}

func (h *Handler) HandleDeleteOilFillLevel(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		handleError(w, err, 500)
		return
	}

	if err := h.store.DeleteOilFillLevel(id); err != nil {
		handleError(w, err, 404)
		return
	}

	http.Redirect(w, r, "/oil", http.StatusFound)
}
