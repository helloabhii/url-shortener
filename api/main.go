package main

import (
	"log"
	"os" // to access env variable
	"url-shortner/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app)
	log.Fatal(app.Listen(os.Getenv("APP_PORT"))) //Create the sever here

}
