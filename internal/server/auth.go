package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"policyAuth/internal/database"

	"github.com/Nerzal/gocloak/v13"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	keycloakClient *gocloak.GoCloak
	realm          string
	clientID       string
	clientSecret   string
)

func InitKeycloak() {
	keycloakClient = gocloak.NewClient(fmt.Sprintf("%v:%v", os.Getenv("KEYCLOAK_URL"), os.Getenv("KEYCLOAK_PORT")))
	realm = os.Getenv("KEYCLOAK_REALM")
	clientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
}

func JWTMiddleware(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Authorization header is required")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		rptResult, err := keycloakClient.RetrospectToken(context.Background(), tokenString, clientID, clientSecret, realm)
		if err != nil || !*rptResult.Active {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		claims := jwt.MapClaims{}
		_, _, err = new(jwt.Parser).ParseUnverified(tokenString, claims)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		c.Locals("user", claims)
		return next(c)
	}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthHandler(c *fiber.Ctx) error {
	var authReq AuthRequest
	if err := c.BodyParser(&authReq); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	token, err := keycloakClient.Login(context.Background(), clientID, clientSecret, realm, authReq.Username, authReq.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("Failed to authenticate user: %v", err))
	}

	// Check if the user exists in the database
	db := database.New().(*database.StructService).DB
	var userID int
	err = db.QueryRow("SELECT user_id FROM pds_users WHERE username = $1", authReq.Username).Scan(&userID)
	if err == sql.ErrNoRows {
		// User does not exist, insert the user
		err = db.QueryRow(
			"INSERT INTO pds_users (username, password) VALUES ($1, $2) RETURNING user_id",
			authReq.Username, authReq.Password,
		).Scan(&userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to add user to the database")
		}

		// Assign default role to the new user (e.g., "user")
		var roleID int
		err = db.QueryRow("SELECT role_id FROM pds_roles WHERE role_name = $1", "user").Scan(&roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to find default role")
		}

		_, err = db.Exec("INSERT INTO pds_user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to assign role to the user")
		}
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error")
	}

	authResp := AuthResponse{Token: token.AccessToken}
	return c.JSON(authResp)
}
