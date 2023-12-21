package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/helloabhii/url-shortner/helpers"

	"github.com/helloabhii/url-shortner/database"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid" //give unique id for every single user
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customshort"`
	Expiry      time.Duration `json:"expiry"`
}
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"customshort"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"_rate_remaining"`
	XRateLimitReset time.Duration `json:"x_rate_limitrest"`
}

func ShortenURL(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error ": "cannot parse JSON"})
	}
	//implement rate limiting
	r2 := database.CreateClient(1)                    //goint to database client
	defer r2.Close()                                  // closing the connection to this database
	val, err := r2.Get(database.Ctx, c.IP()).Result() //check the ip address  // .Result() -> because i want get result back
	if err == redis.Nil {                             //means you didn't find any value in the database
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err() // first time client using this database // Err() -> if any issue occurs
	} else { //if user already found in the database
		val, _ = r2.Get(database.Ctx, c.IP()).Result() // _ -> because here we didn't define the error
		valInt, _ := strconv.Atoi(val)                 //converting it into int
		if valInt <= 0 {                               //if valint show 0 means you can't use this service for now
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate limit reset": limit / time.Nanosecond / time.Minute,
			})
		}

	}

	//check if the input by the user actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "you can't hack the system(: is"})
	}

	//enforce https, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	var id string
	if body.CustomShort == " " {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()
	//checking that value actually (url) already used by someone or not
	val, _ = r.Get(database.Ctx, id).Result() //id could be customshort or that user give
	if val != " " {                           //something is there
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use ",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24 //setting 24 hour
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err() //id -> custom short url or real

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to server",
		})
	}

	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP()) // decrementing the value of rate remaining

	//we again need the value
	val, _ = r2.Get(database.Ctx, c.IP()).Result() //value fetch from database
	resp.XRateRemaining, _ = strconv.Atoi(val)     //

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)
}
