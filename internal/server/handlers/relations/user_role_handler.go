package relations

import (
	"database/sql"

	models "policyAuth/internal/models/relations"

	"github.com/gofiber/fiber/v2"
)

type UserRoleHandler struct {
	DB *sql.DB
}

func (h *UserRoleHandler) userExists(userID int) (bool, error) {
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pds_users WHERE user_id=$1)", userID).Scan(&exists)
	return exists, err
}

func (h *UserRoleHandler) roleExists(roleID int) (bool, error) {
	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pds_roles WHERE role_id=$1)", roleID).Scan(&exists)
	return exists, err
}

func (h *UserRoleHandler) validateUserRole(c *fiber.Ctx, userID, roleID int) error {
	userExists, err := h.userExists(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if !userExists {
		return c.Status(fiber.StatusBadRequest).SendString("User does not exist")
	}

	roleExists, err := h.roleExists(roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if !roleExists {
		return c.Status(fiber.StatusBadRequest).SendString("Role does not exist")
	}

	return nil
}

func (h *UserRoleHandler) GetUserRoles(c *fiber.Ctx) error {
	rows, err := h.DB.Query(`
    SELECT rl.user_id, pu.username , rl.role_id, pr.role_name 
    FROM pds_user_roles rl
    join pds_roles pr on pr.role_id = rl.role_id
    join pds_users pu on pu.user_id = rl.user_id
    `)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rows.Close()

	var userRoles []models.UserRole
	for rows.Next() {
		var userRole models.UserRole
		if err := rows.Scan(&userRole.UserID,&userRole.UserName, &userRole.RoleID,&userRole.RoleName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		userRoles = append(userRoles, userRole)
	}

	return c.JSON(userRoles)
}

func (h *UserRoleHandler) CreateUserRole(c *fiber.Ctx) error {
	var userRole models.UserRole
	if err := c.BodyParser(&userRole); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.validateUserRole(c, userRole.UserID, userRole.RoleID); err != nil {
		return err
	}

	_, err := h.DB.Exec("INSERT INTO pds_user_roles (user_id, role_id) VALUES ($1, $2)", userRole.UserID, userRole.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(userRole)
}

func (h *UserRoleHandler) DeleteUserRole(c *fiber.Ctx) error {
	var userRole models.UserRole
	if err := c.BodyParser(&userRole); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.validateUserRole(c, userRole.UserID, userRole.RoleID); err != nil {
		return err
	}

	_, err := h.DB.Exec("DELETE FROM pds_user_roles WHERE user_id=$1 AND role_id=$2", userRole.UserID, userRole.RoleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
