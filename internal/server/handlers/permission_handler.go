package handlers

import (
    "database/sql"
    "strconv"

    "policyAuth/internal/models"
    "policyAuth/internal/services"

    "github.com/gofiber/fiber/v2"
)

type PermissionHandler struct {
    PermissionService *services.PermissionService
}

func NewPermissionHandler(DB *sql.DB) *PermissionHandler {
    return &PermissionHandler{PermissionService: services.NewPermissionService(DB)}
}

func (h *PermissionHandler) GetPermissions(c *fiber.Ctx) error {
    permissions, err := h.PermissionService.GetPermissions()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(permissions)
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
    var permission models.Permission
    if err := c.BodyParser(&permission); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    createdPermission, err := h.PermissionService.CreatePermission(permission)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(createdPermission)
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

    if err := h.PermissionService.UpdatePermission(id, permission); err != nil {
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

    if err := h.PermissionService.DeletePermission(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

