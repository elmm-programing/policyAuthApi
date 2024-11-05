package handlers

import (
    "database/sql"
    "strconv"

    "policyAuth/internal/models"
    "policyAuth/internal/services"

    "github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
    RoleService *services.RoleService
}

func NewRoleHandler(DB *sql.DB) *RoleHandler {
    return &RoleHandler{RoleService: services.NewRoleService(DB)}
}

func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
    roles, err := h.RoleService.GetRoles()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(roles)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
    var role models.Role
    if err := c.BodyParser(&role); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    createdRole, err := h.RoleService.CreateRole(role)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(createdRole)
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

    if err := h.RoleService.UpdateRole(id, role); err != nil {
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

    if err := h.RoleService.DeleteRole(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

