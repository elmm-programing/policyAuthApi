package relations 

type ResourceRole struct {
  ID         int `json:"id"`
	ResourceID int `json:"resource_id"`
	RoleID     int `json:"role_id"`
}

type RoleResourcePermission struct {
	ID                      int `json:"id"`
	ResourceRoleID int `json:"resource_role_id"`
	PermissionID            int `json:"permission_id"`
}




