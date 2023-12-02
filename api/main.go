package main

import (
	"fmt"
	"log"
	"os" // to access env variable

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
	}
	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT"))) //Create the sever here

}
