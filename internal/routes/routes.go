package routes

import (
	"gin-test/internal/handlers"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(app *fiber.App) {
	app.Get("/ready", handlers.Ready)
	app.Post("/fraud-score", handlers.ApproveTransaction)
	app.Get("/references", handlers.LoadJson)
}
