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

func HandleAddWaterGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/water/add.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleAddWaterPost(w http.ResponseWriter, r *http.Request) {
	year, err := utils.ConvertInt(r.FormValue("year"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volumeWater, err := utils.ConvertFloat(r.FormValue("volumeWater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volumeWastewater, err := utils.ConvertFloat(r.FormValue("volumeWastewater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volumeRainwater, err := utils.ConvertFloat(r.FormValue("volumeRainwater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costsWater, err := utils.ConvertFloat(r.FormValue("costsWater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costsWastewater, err := utils.ConvertFloat(r.FormValue("costsWastewater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costsRainwater, err := utils.ConvertFloat(r.FormValue("costsRainwater"))
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
	fixedPrice, err := utils.ConvertFloat(r.FormValue("fixedPrice"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/water", http.StatusFound)
}

func HandleEditWater(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("edit.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/water/edit.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	entry, err := storage.GetWaterEntry(id)
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

func HandleSaveWater(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	year, err := utils.ConvertInt(r.FormValue("year"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volumeWater, err := utils.ConvertFloat(r.FormValue("volumeWater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volumeWastewater, err := utils.ConvertFloat(r.FormValue("volumeWastewater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volumeRainwater, err := utils.ConvertFloat(r.FormValue("volumeRainwater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costsWater, err := utils.ConvertFloat(r.FormValue("costsWater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costsWastewater, err := utils.ConvertFloat(r.FormValue("costsWastewater"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	costsRainwater, err := utils.ConvertFloat(r.FormValue("costsRainwater"))
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
	fixedPrice, err := utils.ConvertFloat(r.FormValue("fixedPrice"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/water", http.StatusFound)
}

func HandleDeleteWater(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	err = storage.DeleteWater(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/water", http.StatusFound)
}
