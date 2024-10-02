package handlers

import (
	"database/sql"
	"strconv"

	"policyAuth/internal/models"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT user_id, username, password FROM pds_users")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.Password); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		users = append(users, user)
	}

	return c.JSON(users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.DB.QueryRow("INSERT INTO pds_users (username, password) VALUES ($1, $2) RETURNING user_id", user.Username, user.Password).Scan(&user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(user)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err = h.DB.Exec("UPDATE pds_users SET username=$1, password=$2 WHERE user_id=$3", user.Username, user.Password, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	_, err = h.DB.Exec("DELETE FROM pds_users WHERE user_id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
