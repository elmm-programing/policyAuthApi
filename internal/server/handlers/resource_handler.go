package handlers

import (
    "database/sql"
    "strconv"

    "policyAuth/internal/models"
    "policyAuth/internal/services"

    "github.com/gofiber/fiber/v2"
)

type ResourceHandler struct {
    ResourceService *services.ResourceService
}

func NewResourceHandler(DB *sql.DB) *ResourceHandler {
    return &ResourceHandler{ResourceService: services.NewResourceService(DB)}
}

func (h *ResourceHandler) GetResources(c *fiber.Ctx) error {
    resources, err := h.ResourceService.GetResources()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(resources)
}

func (h *ResourceHandler) CreateResource(c *fiber.Ctx) error {
    var resource models.Resource
    if err := c.BodyParser(&resource); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    createdResource, err := h.ResourceService.CreateResource(resource)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(createdResource)
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

    if err := h.ResourceService.UpdateResource(id, resource); err != nil {
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

    if err := h.ResourceService.DeleteResource(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

