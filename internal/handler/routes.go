package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /", h.HandleView)

	mux.HandleFunc("GET /web/", h.HandleFiles)
	mux.HandleFunc("GET /service-worker.js", h.HandleServiceWorker)

	mux.HandleFunc("GET /electricity", h.HandleElectricityView)
	mux.HandleFunc("GET /electricity/add", h.HandleAddElectricityGet)
	mux.HandleFunc("POST /electricity/add", h.HandleAddElectricityPost)
	mux.HandleFunc("GET /electricity/edit/{id}", h.HandleEditElectricity)
	mux.HandleFunc("POST /electricity/save/{id}", h.HandleSaveElectricity)
	mux.HandleFunc("GET /electricity/delete/{id}", h.HandleDeleteElectricity)

	mux.HandleFunc("GET /oil", h.HandleOilView)
	mux.HandleFunc("GET /oil/add", h.HandleAddOilGet)
	mux.HandleFunc("POST /oil/add", h.HandleAddOilPost)
	mux.HandleFunc("GET /oil/edit/{id}", h.HandleEditOil)
	mux.HandleFunc("POST /oil/save/{id}", h.HandleSaveOil)
	mux.HandleFunc("GET /oil/delete/{id}", h.HandleDeleteOil)

	mux.HandleFunc("GET /oil/fill-level/add", h.HandleAddOilFillLevelGet)
	mux.HandleFunc("POST /oil/fill-level/add", h.HandleAddOilFillLevelPost)
	mux.HandleFunc("GET /oil/fill-level/edit/{id}", h.HandleEditOilFillLevel)
	mux.HandleFunc("POST /oil/fill-level/save/{id}", h.HandleSaveOilFillLevel)
	mux.HandleFunc("GET /oil/fill-level/delete/{id}", h.HandleDeleteOilFillLevel)

	mux.HandleFunc("GET /water", h.HandleWaterView)
	mux.HandleFunc("GET /water/add", h.HandleAddWaterGet)
	mux.HandleFunc("POST /water/add", h.HandleAddWaterPost)
	mux.HandleFunc("GET /water/edit/{id}", h.HandleEditWater)
	mux.HandleFunc("POST /water/save/{id}", h.HandleSaveWater)
	mux.HandleFunc("GET /water/delete/{id}", h.HandleDeleteWater)
}
