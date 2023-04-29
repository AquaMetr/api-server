package main

import (
	"github.com/AquaMetr/api-server/internal/api"
	"github.com/AquaMetr/api-server/internal/storage"
	"github.com/labstack/gommon/log"
)

func main() {

	log.Info("Starting AquaMetr api-server")

	store := storage.NewStorage()

	log.Fatal(api.StartApi(store))
}
