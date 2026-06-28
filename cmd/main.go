package main

import (
	"log"
	"net/http"

	"github.com/lukas-arnold/consumption-tracker/internal/configs"
	"github.com/lukas-arnold/consumption-tracker/internal/handler"
	"github.com/lukas-arnold/consumption-tracker/internal/language"
	"github.com/lukas-arnold/consumption-tracker/internal/storage"
)

func main() {
	if err := language.LoadLanguages(); err != nil {
		log.Fatal(err)
	}

	store := storage.New(configs.GetStorageFile())
	h := handler.New(store)
	server := createServer(h)

	log.Printf("Consumption Tracker running on %s", configs.GetPort())
	log.Fatal(http.ListenAndServe(configs.GetPort(), server))
}
