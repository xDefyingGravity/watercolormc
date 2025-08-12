package app

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"os/signal"
	"syscall"
	"watercolormc/internal/app/channels"
	"watercolormc/internal/app/middleware"
	"watercolormc/internal/app/routes"
)

func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		BodyLimit:             10 * 1024 * 1024 * 1024,
	})

	middleware.Setup(app)
	routes.Setup(app)
	channels.Init(app)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		channels.Cleanup()
		os.Exit(0)
	}()

	return app
}
