package main

import (
	"log"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/handler"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
)

func main() {
	err := language.LoadLanguages()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /service-worker", handler.HandleServiceWorker)
	mux.HandleFunc("GET /web/", handler.HandleFiles)

	mux.HandleFunc("GET /", handler.HandleView)

	mux.HandleFunc("GET /electricity", handler.HandleElectricityView)
	mux.HandleFunc("GET /oil", handler.HandleOilView)
	mux.HandleFunc("GET /water", handler.HandleWaterView)

	mux.HandleFunc("GET /electricity/add", handler.HandleAddElectricityGet)
	mux.HandleFunc("POST /electricity/add", handler.HandleAddElectricityPost)
	mux.HandleFunc("GET /electricity/edit/{id}", handler.HandleEditElectricity)
	mux.HandleFunc("POST /electricity/save/{id}", handler.HandleSaveElectricity)
	mux.HandleFunc("GET /electricity/delete/{id}", handler.HandleDeleteElectricity)

	mux.HandleFunc("GET /oil/add", handler.HandleAddOilGet)
	mux.HandleFunc("POST /oil/add", handler.HandleAddOilPost)
	mux.HandleFunc("GET /oil/edit/{id}", handler.HandleEditOil)
	mux.HandleFunc("POST /oil/save/{id}", handler.HandleSaveOil)
	mux.HandleFunc("GET /oil/delete/{id}", handler.HandleDeleteOil)
	mux.HandleFunc("GET /oil-fill-level/add", handler.HandleAddOilFillLevelGet)
	mux.HandleFunc("POST /oil-fill-level/add", handler.HandleAddOilFillLevelPost)
	mux.HandleFunc("GET /oil-fill-level/edit/{id}", handler.HandleEditOilFillLevel)
	mux.HandleFunc("POST /oil-fill-level/save/{id}", handler.HandleSaveOilFillLevel)
	mux.HandleFunc("GET /oil-fill-level/delete/{id}", handler.HandleDeleteOilFillLevel)

	mux.HandleFunc("GET /water/add", handler.HandleAddWaterGet)
	mux.HandleFunc("POST /water/add", handler.HandleAddWaterPost)
	mux.HandleFunc("GET /water/edit/{id}", handler.HandleEditWater)
	mux.HandleFunc("POST /water/save/{id}", handler.HandleSaveWater)
	mux.HandleFunc("GET /water/delete/{id}", handler.HandleDeleteWater)

	log.Printf("Consumption Tracker running on %s", configs.GetPort())
	log.Fatal(http.ListenAndServe(configs.GetPort(), mux))
}
