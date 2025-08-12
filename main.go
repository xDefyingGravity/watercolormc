package main

import (
	"go.uber.org/zap"
	"watercolormc/internal"
	"watercolormc/internal/app"
	"watercolormc/internal/app/channels"
	"watercolormc/internal/database"
	"watercolormc/internal/logger"
)

func main() {
	log := logger.Init()

	if err := database.Init(); err != nil {
		log.Fatal(err.Error())
	}

	if err := database.SetupSchema(); err != nil {
		log.Fatal(err.Error())
	}

	server := app.Init()
	defer channels.Cleanup()

	log.Info("starting server", zap.String("port", internal.PORT))
	if err := server.Listen(internal.PORT); err != nil {
		log.Fatal(err.Error())
	}
}
