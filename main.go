package main

import (
	"fmt"
	"gin-test/internal/ivf"
	"gin-test/internal/routes"
	"gin-test/internal/store"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	fmt.Println("Building IVF index...")
	ivf.Clusters = ivf.BuildIVFIndexStreamed(store.References, 512, 1)
	fmt.Println("Initializing Server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	app := fiber.New()
	routes.RegisterRoutes(app)
	app.Listen(":"+port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})

}

func init() {
	fmt.Println("Loading references")
	store.LoadReferencesGziped(10000)
	store.LoadNormalizer()
}
