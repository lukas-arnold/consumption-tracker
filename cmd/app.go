package main

import (
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/handler"
)

func createServer(
	h *handler.Handler,
) http.Handler {

	mux := http.NewServeMux()

	handler.RegisterRoutes(
		mux,
		h,
	)

	return mux
}
