package routes

import (
	"example/url-shortner/helpers"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/helloabhii/url-shortner/database"
	"github.com/helloabhii/url-shortner/helpers"
)

type request struct {
	URL         string        `json: "url"`
	CustomShort string        `json: "customshort"`
	Expiry      time.Duration `json: "expiry"`
}
type response struct {
	URL             string        `json: "url"`
	CustomShort     string        `json: "customshort"`
	Expiry          time.Duration `json: "expiry"`
	XRateRemaining  int           `json: "x_rate_remaining`
	XRateLimitReset time.Duration `json:	"x_rate_limitrest"`
}

func shorternURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error ": "cannot parse JSON"})
	}
	//implement rate limiting
	r2 := database.CreateClient(1)
	defer r2.close()
	val, err := r2.Get(databaseCtx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.IP().Result())
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate limit reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	//check if the input by the user actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest.JSON(fiber.Map{""}))

	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber / StatusServiceUnavailable).JSON(fiber.Map{"error": "you can't hack the system(: is"})
	}

	//enforce https, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	r2.Decr(database.Ctx, c.IP())

}
