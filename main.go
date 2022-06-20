package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/woonmapao/go-fiber-mongo/controllers"
)

func main() {
	if err := controllers.Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Get("/employee", controllers.GetEmployee)

	app.Post("/employee", controllers.CreateEmployee)

	app.Put("/employee/:id", controllers.EditEmployee)

	app.Delete("/employee/:id", controllers.DeleteEmployee)

	log.Fatal(app.Listen(":3000"))
}
