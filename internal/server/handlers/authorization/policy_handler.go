package authorization

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ResourceDetailsHandler struct {
	DB *sql.DB
}

func (h *ResourceDetailsHandler) userExists(username string) (bool, int, error) {
	var userID int
	err := h.DB.QueryRow("SELECT user_id FROM pds_users WHERE username=$1", username).Scan(&userID)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, userID, nil
}

func (h *ResourceDetailsHandler) GetRolesAndPermissionsForResource(c *fiber.Ctx) error {
	username := c.Params("username")
  userExists, userID, err := h.userExists(username)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if !userExists {
		return c.Status(fiber.StatusBadRequest).SendString("User does not exist")
	}
  userRolesRows, err := h.DB.Query(`
		SELECT r.role_name
		FROM pds_roles r
		JOIN pds_user_roles ur ON r.role_id = ur.role_id
		WHERE ur.user_id = $1`, userID)
if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
  var userRoles []string
	for userRolesRows.Next() {
		var roleName string
		if err := userRolesRows.Scan(&roleName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		userRoles = append(userRoles, roleName)
	}
if len(userRoles) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("User has no roles")
	}
query := `
		SELECT DISTINCT res.resource_name
		FROM pds_resources res
		JOIN pds_resource_role rr ON res.resource_id = rr.resource_id
		JOIN pds_roles r ON rr.role_id = r.role_id
		WHERE r.role_name = ANY($1)`
	resourceRows, err := h.DB.Query(query, "{"+strings.Join(userRoles, ",")+"}")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer resourceRows.Close()

	var resources []string
	for resourceRows.Next() {
		var resourceName string
		if err := resourceRows.Scan(&resourceName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		resources = append(resources, resourceName)
	}



	// Create the response structure
  response := make(map[string][]string)
  for _,resource := range resources {
    response[resource] = userRoles

  }

	return c.JSON(response)
  // return c.Status(fiber.StatusNotFound).SendString("Not found")
}
