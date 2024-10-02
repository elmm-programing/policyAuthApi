package relations 

import (
	"database/sql"
	 "policyAuth/internal/models/relations"

	"github.com/gofiber/fiber/v2"
)

type ResourceRoleHandler struct {
	DB *sql.DB
}

func (h *ResourceRoleHandler) GetResourceRoles(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, resource_id, role_id FROM pds_resource_role")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var resourceRoles []relations.ResourceRole
	for rows.Next() {
		var resourceRole relations.ResourceRole
		if err := rows.Scan(&resourceRole.ID, &resourceRole.ResourceID, &resourceRole.RoleID); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		resourceRoles = append(resourceRoles, resourceRole)
	}

	return c.JSON(resourceRoles)
}

func (h *ResourceRoleHandler) CreateResourceRole(c *fiber.Ctx) error {
	var resourceRole relations.ResourceRole
	if err := c.BodyParser(&resourceRole); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err := h.DB.Exec("INSERT INTO pds_resource_role (resource_id, role_id) VALUES ($1, $2)", resourceRole.ResourceID, resourceRole.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(resourceRole)
}

func (h *ResourceRoleHandler) DeleteResourceRole(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.DB.Exec("DELETE FROM pds_resource_role WHERE id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type RoleResourcePermissionHandler struct {
	DB *sql.DB
}

func (h *RoleResourcePermissionHandler) GetRoleResourcePermissions(c *fiber.Ctx) error {
	rows, err := h.DB.Query("SELECT id, resource_role_id, permission_id FROM pds_role_resource_permissions")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var roleResourcePermissions []relations.RoleResourcePermission
	for rows.Next() {
		var roleResourcePermission relations.RoleResourcePermission
		if err := rows.Scan(&roleResourcePermission.ID, &roleResourcePermission.ResourceRoleID, &roleResourcePermission.PermissionID); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		roleResourcePermissions = append(roleResourcePermissions, roleResourcePermission)
	}

	return c.JSON(roleResourcePermissions)
}

func (h *RoleResourcePermissionHandler) CreateRoleResourcePermission(c *fiber.Ctx) error {
	var roleResourcePermission relations.RoleResourcePermission
	if err := c.BodyParser(&roleResourcePermission); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	_, err := h.DB.Exec("INSERT INTO pds_role_resource_permissions (resource_role_id, permission_id) VALUES ($1, $2)", roleResourcePermission.ResourceRoleID, roleResourcePermission.PermissionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(roleResourcePermission)
}

func (h *RoleResourcePermissionHandler) DeleteRoleResourcePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := h.DB.Exec("DELETE FROM pds_role_resource_permissions WHERE id=$1", id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}


