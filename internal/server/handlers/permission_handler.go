package handlers

import (
	"database/sql"
	"strconv"

	"policyAuth/internal/models"

	"github.com/gofiber/fiber/v2"
)

type PermissionHandler struct {
	DB *sql.DB
}

func (h *PermissionHandler) GetPermissions(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT permission_id, permission_name FROM pds_permissions")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		if err := rows.Scan(&permission.PermissionID, &permission.PermissionName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		permissions = append(permissions, permission)
	}

	return c.JSON(permissions)
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	var permission models.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.DB.QueryRow("INSERT INTO pds_permissions (permission_name) VALUES ($1) RETURNING permission_id", permission.PermissionName).Scan(&permission.PermissionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(permission)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid permission ID")
	}

	var permission models.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err = h.DB.Exec("UPDATE pds_permissions SET permission_name=$1 WHERE permission_id=$2", permission.PermissionName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(permission)
}

func (h *PermissionHandler) DeletePermission(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid permission ID")
	}

	_, err = h.DB.Exec("DELETE FROM pds_permissions WHERE permission_id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
