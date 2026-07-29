package handlers

import "github.com/gofiber/fiber/v2"

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeError(c *fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(ErrorResponse{Error: err.Error()})
}
