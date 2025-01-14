package relations

import (
	"database/sql"
	models "policyAuth/internal/models/relations"
	services "policyAuth/internal/services/relations"

	"github.com/gofiber/fiber/v2"
)

type UserRoleHandler struct {
    UserRoleService services.UserRoleService
}

func NewUserRoleHandler(DB *sql.DB) *UserRoleHandler  {
  return &UserRoleHandler{services.UserRoleService{DB: DB}}
  
}

func (h *UserRoleHandler) validateUserRoleById(c *fiber.Ctx, userID, roleID int) error {
    userExists, err := h.UserRoleService.UserExistsById(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    if !userExists {
        return c.Status(fiber.StatusBadRequest).SendString("User does not exist")
    }

    roleExists, err := h.UserRoleService.RoleExistsById(roleID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    if !roleExists {
        return c.Status(fiber.StatusBadRequest).SendString("Role does not exist")
    }

    return nil
}

func (h *UserRoleHandler) validateUserRoleByName(c *fiber.Ctx, userName, roleName string) error {
    userExists, err := h.UserRoleService.UserExistsByName(userName)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    if !userExists {
        return c.Status(fiber.StatusBadRequest).SendString("User does not exist")
    }

    roleExists, err := h.UserRoleService.RoleExistsByName(roleName)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }
    if !roleExists {
        return c.Status(fiber.StatusBadRequest).SendString("Role does not exist")
    }

    return nil
}

func (h *UserRoleHandler) GetUserRoles(c *fiber.Ctx) error {
    userRoles, err := h.UserRoleService.GetUserRoles()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(userRoles)
}

func (h *UserRoleHandler) CreateUserRoleById(c *fiber.Ctx) error {
    var userRole models.UserRole
    if err := c.BodyParser(&userRole); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.validateUserRoleById(c, userRole.UserID, userRole.RoleID); err != nil {
        return err
    }

    if err := h.UserRoleService.CreateUserRole(userRole); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(userRole)
}

func (h *UserRoleHandler) DeleteUserRoleById(c *fiber.Ctx) error {
    var userRole models.UserRole
    if err := c.BodyParser(&userRole); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.validateUserRoleById(c, userRole.UserID, userRole.RoleID); err != nil {
        return err
    }

    if err := h.UserRoleService.DeleteUserRole(userRole); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserRoleHandler) CreateUserRoleByName(c *fiber.Ctx) error {
    var userRole models.UserRole
    if err := c.BodyParser(&userRole); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.validateUserRoleByName(c, userRole.UserName, userRole.RoleName); err != nil {
        return err
    }

    if err := h.UserRoleService.CreateUserRole(userRole); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(userRole)
}

func (h *UserRoleHandler) DeleteUserRoleByName(c *fiber.Ctx) error {
    var userRole models.UserRole
    if err := c.BodyParser(&userRole); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.validateUserRoleById(c, userRole.UserID, userRole.RoleID); err != nil {
        return err
    }

    if err := h.UserRoleService.DeleteUserRole(userRole); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

