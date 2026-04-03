package main

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CreateRequest struct {
	Name string `json:"name"`
	A    string `json:"a"`
	B    string `json:"b"`
}

type UpdateRequest struct {
	Name          string `json:"name"`
	A             string `json:"a"`
	B             string `json:"b"`
	ConflictToken string `json:"conflict_token"`
}

func registerRoutes(app *fiber.App, strategy ConcurrencyStrategy) {
	app.Post("/test", func(c fiber.Ctx) error {
		var req CreateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		id, token, err := strategy.Insert(req.Name, req.A, req.B)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":             id,
			"conflict_token": token,
		})
	})

	app.Get("/test/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		record, err := strategy.GetByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(record)
	})

	app.Put("/test/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		var req UpdateRequest
		if err := c.Bind().JSON(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		if req.ConflictToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "conflict_token is required",
			})
		}

		newToken, err := strategy.Update(id, req.Name, req.A, req.B, req.ConflictToken)
		if err != nil {
			if errors.Is(err, ErrConflict) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Conflict: record was modified by another transaction",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":        "Updated successfully",
			"conflict_token": newToken,
		})
	})
}
