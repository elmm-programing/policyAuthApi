package server

import (
	"fmt"
	"os"
	"strconv"

	"policyAuth/internal/database"
	"policyAuth/internal/helpers"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   *database.DatabaseService
  helpers  helpers.Helpers
	app  *fiber.App
}

func NewServer() *Server {
	InitKeycloak()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbService := database.New()
	database.InitSchema(dbService.Instance) // Initialize the database schema

	app := fiber.New()

	server := &Server{
		port: port,
		db:   dbService,
		app:  app,
	}
	server.RegisterRoutes()

	return server
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%d", s.port))
}
