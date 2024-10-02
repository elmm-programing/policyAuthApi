package relations 

type ResourceRole struct {
  ID         int `json:"id"`
	ResourceID int `json:"resource_id"`
	RoleID     int `json:"role_id"`
}

type RoleResourcePermission struct {
	ID                      int `json:"id"`
	ResourceRolePermissionID int `json:"resource_role_permission_id"`
	PermissionID            int `json:"permission_id"`
}

