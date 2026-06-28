package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

type Handler struct {
	store *storage.Storage
}

func New(
	store *storage.Storage,
) *Handler {
	return &Handler{
		store: store,
	}
}

func handleError(
	w http.ResponseWriter,
	err error,
	statusCode int,
) {
	log.Printf(
		"HTTP %d: %v",
		statusCode,
		err,
	)

	http.Error(
		w,
		http.StatusText(statusCode),
		statusCode,
	)
}

func (h *Handler) renderTemplate(
	w http.ResponseWriter,
	file string,
	data any,
) {
	tmpl := template.Must(
		template.New("base.html").
			Funcs(getTemplateFuncs()).
			ParseFS(
				configs.GetWebFiles(),
				"templates/base.html",
				file,
			),
	)

	if err := tmpl.Execute(
		w,
		data,
	); err != nil {
		handleError(
			w,
			err,
			500,
		)
	}
}

func formatFloat(
	value float64,
	digits int,
) string {

	formatted :=
		fmt.Sprintf(
			"%.*f",
			digits,
			value,
		)

	return strings.ReplaceAll(
		formatted,
		".",
		",",
	)
}

func formatInt(
	value float64,
) string {

	formatted :=
		fmt.Sprintf(
			"%.0f",
			value,
		)

	return strings.ReplaceAll(
		formatted,
		".",
		",",
	)
}

func formatEuro(
	value float64,
) string {

	if value == 0 {
		return "-"
	}

	return formatFloat(
		value,
		2,
	) + " €"
}

func formatPrice(
	costs float64,
	units float64,
	unit string,
) string {

	if units <= 0 {
		return "-"
	}

	return formatFloat(
		costs/units,
		3,
	) + " €/" + unit
}

func calculateMonthlyPeriod(
	payments float64,
	days float64,
) float64 {

	return payments / days * 30.4375
}

func formatMonthly(
	payments float64,
	timeFrom string,
	timeTo string,
) string {

	if payments == 0 {
		return "-"
	}

	monthly :=
		payments / 12

	if timeFrom == "" ||
		timeTo == "" {

		return formatEuro(monthly)
	}

	from, errFrom :=
		utils.ConvertTime(timeFrom)

	to, errTo :=
		utils.ConvertTime(timeTo)

	if errFrom != nil ||
		errTo != nil {

		return formatEuro(monthly)
	}

	days :=
		to.Sub(from).
			Hours() / 24

	if days <= 0 ||
		(days >= 364 && days <= 366) {

		return formatEuro(monthly)
	}

	return formatEuro(
		calculateMonthlyPeriod(
			payments,
			days,
		),
	)
}

func formatDifference(
	payments float64,
	costs float64,
) string {

	difference :=
		payments - costs

	if difference == 0 {
		return "-"
	}

	return formatEuro(
		difference,
	)
}

func formatTotalCosts(
	costsWater float64,
	costsWastewater float64,
	costsRainwater float64,
	fixedPrice float64,
) string {

	total :=
		costsWater +
			costsWastewater +
			costsRainwater +
			fixedPrice

	if total == 0 {
		return "-"
	}

	return formatEuro(
		total,
	)
}

func add(
	values ...float64,
) float64 {

	total := 0.0

	for _, value := range values {
		total += value
	}

	return total
}

func getTemplateFuncs() template.FuncMap {
	return template.FuncMap{

		"T": func(
			key string,
		) string {
			return language.T(
				configs.GetLanguage(),
				key,
			)
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

func (h *Handler) HandleView(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.HandleElectricityView(
		w,
		r,
	)
}
