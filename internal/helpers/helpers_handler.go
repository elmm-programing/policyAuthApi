package helpers 

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Helpers) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}

func (s *Helpers) HealthHandler(c *fiber.Ctx) error {
	return c.JSON(s.DB.Health())
}
