package handler

import (
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
)

func (h *Handler) HandleFiles(
	w http.ResponseWriter,
	r *http.Request,
) {
	http.StripPrefix(
		"/web/",
		http.FileServerFS(
			configs.GetWebFiles(),
		),
	).ServeHTTP(
		w,
		r,
	)
}

func (h *Handler) HandleServiceWorker(
	w http.ResponseWriter,
	r *http.Request,
) {

	http.ServeFileFS(
		w,
		r,
		configs.GetWebFiles(),
		"service-worker.js",
	)
}
