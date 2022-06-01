package main

import (
	"database/sql"
	"fmt"
	"link-shortener/models"
	"mime"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	_ = mime.AddExtensionType(".js", "text/javascript")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	{
		if err != nil {
			panic(err)
		}
		// close database
		defer db.Close()

		// check connection
		err = db.Ping()
		if err != nil {
			panic(err)
		}
		fmt.Println("Postgres is connected!")
	}

	app := fiber.New()

	linkService := models.NewLinkService(db)

	app.Static("/", "./fe")

	app.Post("/shorten", func(c *fiber.Ctx) error {
		return c.JSON(linkService.CreateShortLink(c))
	})

	app.Get("/:hash", func(c *fiber.Ctx) error {
		fullLink := linkService.GetFullLink(c.Params("hash"))
		return c.Redirect(fullLink)
	})

	app.Listen(fmt.Sprintf("localhost:%v", os.Getenv("PORT")))
}
