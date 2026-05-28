package main

import (
	"gin-test/internal/routes"
	"gin-test/internal/store"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	routes.RegisterRoutes(app)
	app.Listen(":8090", fiber.ListenConfig{
		DisableStartupMessage: true,
	})

}

func init() {
	store.LoadReferencesGziped(300000)
	store.LoadNormalizer()
}
