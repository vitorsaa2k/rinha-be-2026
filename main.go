package main

import (
	"gin-test/internal/handlers"
	"gin-test/internal/store"
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	if err := fasthttp.ListenAndServe(":"+port, handlers.FastHTTPHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
	/* app := fiber.New()
	routes.RegisterRoutes(app)
	app.Listen(":9999", fiber.ListenConfig{
		DisableStartupMessage: true,
	}) */

}

func init() {
	store.LoadReferencesGziped(10000)
	store.LoadNormalizer()
}
