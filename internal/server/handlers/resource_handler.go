package handlers

import (
	"database/sql"
	"strconv"

	"policyAuth/internal/models"

	"github.com/gofiber/fiber/v2"
)

type ResourceHandler struct {
	DB *sql.DB
}

func (h *ResourceHandler) GetResources(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT resource_id, resource_name FROM pds_resources")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var resources []models.Resource
	for rows.Next() {
		var resource models.Resource
		if err := rows.Scan(&resource.ResourceID, &resource.ResourceName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		resources = append(resources, resource)
	}

	return c.JSON(resources)
}

func (h *ResourceHandler) CreateResource(c *fiber.Ctx) error {
	var resource models.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := h.DB.QueryRow("INSERT INTO pds_resources (resource_name) VALUES ($1) RETURNING resource_id", resource.ResourceName).Scan(&resource.ResourceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(resource)
}

func (h *ResourceHandler) UpdateResource(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid resource ID")
	}

	var resource models.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err = h.DB.Exec("UPDATE pds_resources SET resource_name=$1 WHERE resource_id=$2", resource.ResourceName, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(resource)
}

func (h *ResourceHandler) DeleteResource(c *fiber.Ctx) error {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid resource ID")
	}

	_, err = h.DB.Exec("DELETE FROM pds_resources WHERE resource_id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
