package routes

import (
	"fmt"

	database "github.com/helloabhii/url-shortner/api/database"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {

	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	oriUrl, err := r.Get(database.Ctx, url).Result()
	fmt.Println(oriUrl)
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "short not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to thr DB ",
		})

	}
	rInr := database.CreateClient(1)
	defer rInr.Close()

	return c.Redirect(oriUrl, 301)
}
