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

func HandleAddOilGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/oil/add.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleAddOilPost(w http.ResponseWriter, r *http.Request) {
	volume, err := utils.ConvertFloat(r.FormValue("volume"))
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
	err = storage.AddOil(models.OilInput{
		Date:     r.FormValue("date"),
		Volume:   volume,
		Costs:    costs,
		Retailer: r.FormValue("retailer"),
		Note:     r.FormValue("note"),
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleEditOil(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("edit.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/oil/edit.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	entry, err := storage.GetOilEntry(id)
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

func HandleSaveOil(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	volume, err := utils.ConvertFloat(r.FormValue("volume"))
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
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleDeleteOil(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	err = storage.DeleteOil(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleAddOilFillLevelGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("add_fill_level.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/oil/add_fill_level.html"),
	)
	err := tmpl.Execute(w, nil)
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
	}
}

func HandleAddOilFillLevelPost(w http.ResponseWriter, r *http.Request) {
	level, err := utils.ConvertFloat(r.FormValue("level"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	err = storage.AddOilFillLevel(models.OilFillLevelInput{
		Date:  r.FormValue("date"),
		Level: level,
	})
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleEditOilFillLevel(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.New("edit_fill_level.html").Funcs(template.FuncMap{
			"T": func(key string) string {
				return language.T(configs.GetLanguage(), key)
			},
		}).ParseFS(configs.GetWebFiles(), "templates/oil/edit_fill_level.html"),
	)
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	entry, err := storage.GetOilFillLevel(id)
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

func HandleSaveOilFillLevel(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	level, err := utils.ConvertFloat(r.FormValue("level"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
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
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}

func HandleDeleteOilFillLevel(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ConvertId(r.PathValue("id"))
	if err != nil {
		errorHandling(w, 500)
		log.Print(err)
		return
	}
	err = storage.DeleteOilFillLevel(id)
	if err != nil {
		errorHandling(w, 404)
		log.Print(err)
		return
	}
	http.Redirect(w, r, "/oil", http.StatusFound)
}
