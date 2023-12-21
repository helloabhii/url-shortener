package routes // same name of packages help us to esay import function packages from other files

import (
	"fmt"

	database "github.com/helloabhii/url-shortner/database"

	"github.com/go-redis/redis/v8" //redis use here  because to check the url that is shortern or original
	"github.com/gofiber/fiber/v2"
)

func ResolveURL(c *fiber.Ctx) error {

	url := c.Params("url")
	r := database.CreateClient(0) //database number -> 0 //database.go file check
	defer r.Close()

	val, err := r.Get(database.Ctx, url).Result() //check the database that url exist or not
	fmt.Println(val)
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "shorten url not found in the database",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot connect to thr DB ",
		})

	}
	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "Couter")
	return c.Redirect(val, 301) //everthing went well then connect to the user
}

//this file contains original url that you give - when you put the url that is shorterned then you will get the link that is shortened by this
// redis is a key value pair database
