package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}

func (s *Server) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
