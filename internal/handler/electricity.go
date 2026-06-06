package handler

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
	"github.com/lukas-arnold/consumption-tracker/internal/models"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
	"github.com/lukas-arnold/consumption-tracker/internal/utils"
)

func HandleAddElectricityGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/electricity/add.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleAddElectricityPost(w http.ResponseWriter, r *http.Request) {
	consumption, err := utils.ConvertFloat(r.FormValue("consumption"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costs, err := utils.ConvertFloat(r.FormValue("costs"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	payments, err := utils.ConvertFloat(r.FormValue("payments"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func HandleEditElectricity(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("edit.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/electricity/edit.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	entry, err := storage.GetElectricity(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}

	err = tmpl.Execute(w, entry)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleSaveElectricity(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	consumption, err := utils.ConvertFloat(r.FormValue("consumption"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costs, err := utils.ConvertFloat(r.FormValue("costs"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	payments, err := utils.ConvertFloat(r.FormValue("payments"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/electricity", http.StatusFound)
}

func HandleDeleteElectricity(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	err = storage.DeleteElectricity(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/electricity", http.StatusFound)
}
