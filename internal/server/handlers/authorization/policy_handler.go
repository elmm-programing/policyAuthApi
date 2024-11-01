package authorization

import (
	"database/sql"
	"strings"

	models "policyAuth/internal/models"
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
	resourceQuery := `
		SELECT DISTINCT res.resource_id , res.resource_name,rr.id
		FROM pds_resources res
		JOIN pds_resource_role rr ON res.resource_id = rr.resource_id
		JOIN pds_roles r ON rr.role_id = r.role_id
		WHERE r.role_name = ANY($1)`
	resourceRows, err := h.DB.Query(resourceQuery, "{"+strings.Join(userRoles, ",")+"}")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer resourceRows.Close()

	var resources []models.ResourceWithRoleRelation
	for resourceRows.Next() {
  var resource models.ResourceWithRoleRelation
		if err := resourceRows.Scan(&resource.ResourceID,&resource.ResourceName,&resource.ResourceRoleId); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		resources = append(resources, resource)
	}
	permissionsOfResourceQuery := `
		SELECT DISTINCT pp.permission_name
		FROM pds_permissions pp 
		JOIN pds_resource_permission prp ON prp.permission_id = pp.permission_id
		JOIN pds_role_resource_permissions rrp ON rrp.permission_id = pp.permission_id
		where prp.resource_id = $1 and rrp.resource_role_id = $2`
  rolesOfResourceQuery := `
		SELECT DISTINCT pr.role_name 
		FROM pds_resource_role prr 
		JOIN pds_roles pr ON pr.role_id = prr.role_id 
		where prr.resource_id = $1`
  for i, v := range resources {
 //Roles   
rolesRows, err := h.DB.Query(rolesOfResourceQuery, v.ResourceID)
    if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer rolesRows.Close()

	for rolesRows.Next() {
		var roleName string
		if err := rolesRows.Scan(&roleName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
      resources[i].Roles = append(resources[i].Roles, roleName)
	}


  //Permissions
	permissionsRows, err := h.DB.Query(permissionsOfResourceQuery, v.ResourceID,v.ResourceRoleId)
    if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer permissionsRows.Close()

	for permissionsRows.Next() {
		var permissionsName string
		if err := permissionsRows.Scan(&permissionsName); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
      resources[i].Permissions = append(resources[i].Permissions, permissionsName)
	}

  }
	
	// Create the response structure
	response := map[string]interface{}{}
	for _, resource := range resources {
		response[resource.ResourceName] = map[string]interface{}{
			"roles": resource.Roles,
      "permissions": resource.Permissions,
		}
	}

	return c.JSON(response)
}
