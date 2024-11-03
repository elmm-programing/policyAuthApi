package relations

import (
    "database/sql"
    "policyAuth/internal/models/relations"
    services "policyAuth/internal/services/relations"

    "github.com/gofiber/fiber/v2"
)

type RoleResourcePermissionHandler struct {
    RoleResourcePermissionService services.RoleResourcePermissionService
}

func NewRoleResourcePermissionHandler(DB *sql.DB) *RoleResourcePermissionHandler {
    return &RoleResourcePermissionHandler{RoleResourcePermissionService: *services.NewRoleResourcePermissionService(DB)}
}

func (h *RoleResourcePermissionHandler) GetRoleResourcePermissions(c *fiber.Ctx) error {
    roleResourcePermissions, err := h.RoleResourcePermissionService.GetRoleResourcePermissions()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(roleResourcePermissions)
}

func (h *RoleResourcePermissionHandler) CreateRoleResourcePermission(c *fiber.Ctx) error {
    var roleResourcePermission relations.RoleResourcePermission
    if err := c.BodyParser(&roleResourcePermission); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.RoleResourcePermissionService.CreateRoleResourcePermission(roleResourcePermission); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(roleResourcePermission)
}

func (h *RoleResourcePermissionHandler) DeleteRoleResourcePermission(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := h.RoleResourcePermissionService.DeleteRoleResourcePermission(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

