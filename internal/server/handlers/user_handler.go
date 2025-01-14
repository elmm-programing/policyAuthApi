package handlers

import (
    "database/sql"
    "strconv"

    "policyAuth/internal/models"
    "policyAuth/internal/services"

    "github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    UserService *services.UserService
}

func NewUserHandler(DB *sql.DB) *UserHandler {
    return &UserHandler{UserService: services.NewUserService(DB)}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    users, err := h.UserService.GetUsers()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    createdUser, err := h.UserService.CreateUser(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(createdUser)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
    idStr := c.Query("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
    }

    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }

    if err := h.UserService.UpdateUser(id, user); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.JSON(user)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
    idStr := c.Query("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
    }

    if err := h.UserService.DeleteUser(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
    }

    return c.SendStatus(fiber.StatusNoContent)
}

