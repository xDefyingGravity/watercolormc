package routes

import (
	"github.com/gofiber/fiber/v2"
	"watercolormc/internal/app/routes/api"
)

func Setup(app *fiber.App) {
	api.RegisterApiRoutes(app)
}
