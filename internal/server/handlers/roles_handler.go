package handlers

import (
	"database/sql"
	"strconv"

	"policyAuth/internal/models"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	DB *sql.DB
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT role_id, role_name FROM pds_roles")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.RoleID, &role.RoleName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		roles = append(roles, role)
	}

	return c.JSON(roles)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.DB.QueryRow("INSERT INTO pds_roles (role_name) VALUES ($1) RETURNING role_id", role.RoleName).Scan(&role.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(role)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid role ID")
	}

	var role models.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err = h.DB.Exec("UPDATE pds_roles SET role_name=$1 WHERE role_id=$2", role.RoleName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(role)
}

func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid role ID")
	}

	_, err = h.DB.Exec("DELETE FROM pds_roles WHERE role_id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
