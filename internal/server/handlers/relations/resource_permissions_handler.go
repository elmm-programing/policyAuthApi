package relations


import (
    "database/sql"
    "policyAuth/internal/models/relations"
    services "policyAuth/internal/services/relations"

    "github.com/gofiber/fiber/v2"
)

// ResourcePermissionHandler handles resource-permission-related requests
type ResourcePermissionHandler struct {
    ResourcePermissionService services.ResourcePermissionService
}

func NewResourcePermissionHandler(DB *sql.DB) *ResourcePermissionHandler {
    return &ResourcePermissionHandler{ResourcePermissionService: *services.NewResourcePermissionService(DB)}
}

func (h *ResourcePermissionHandler) GetResourcePermissions(c *fiber.Ctx) error {
    resourcePermissions, err := h.ResourcePermissionService.GetResourcePermissions()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(resourcePermissions)
}

func (h *ResourcePermissionHandler) CreateResourcePermission(c *fiber.Ctx) error {
    var resourcePermission relations.ResourcePermission
    if err := c.BodyParser(&resourcePermission); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.ResourcePermissionService.CreateResourcePermission(resourcePermission); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(resourcePermission)
}

func (h *ResourcePermissionHandler) DeleteResourcePermission(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := h.ResourcePermissionService.DeleteResourcePermission(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

