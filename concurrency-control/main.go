package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	initDB()

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Fiber v3 is running!")
	})

	app.Post("/echo", func(c fiber.Ctx) error {
		var payload map[string]interface{}
		if err := c.Bind().JSON(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Failed to parse JSON body",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Data received",
			"data":    payload,
		})
	})

	log.Fatal(app.Listen(":8080"))
}
