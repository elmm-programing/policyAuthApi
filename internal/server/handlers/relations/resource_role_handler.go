package relations

import (
	"database/sql"

	"policyAuth/internal/models/relations"

  services "policyAuth/internal/services/relations"
	"github.com/gofiber/fiber/v2"
)


type ResourceRoleHandler struct {
    ResourceRoleService services.ResourceRoleService
}

func NewResourceRoleHandler(DB *sql.DB) *ResourceRoleHandler {
    return &ResourceRoleHandler{ResourceRoleService: *services.NewResourceRoleService(DB)}
}

func (h *ResourceRoleHandler) GetResourceRoles(c *fiber.Ctx) error {
    resourceRoles, err := h.ResourceRoleService.GetResourceRoles()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(resourceRoles)
}

func (h *ResourceRoleHandler) CreateResourceRole(c *fiber.Ctx) error {
    var resourceRole relations.ResourceRole
    if err := c.BodyParser(&resourceRole); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.ResourceRoleService.CreateResourceRole(resourceRole); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(resourceRole)
}

func (h *ResourceRoleHandler) DeleteResourceRole(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := h.ResourceRoleService.DeleteResourceRole(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

