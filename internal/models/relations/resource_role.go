package relations 

type ResourceRole struct {
  ID         int `json:"id"`
	ResourceID int `json:"resource_id"`
	ResourceName string `json:"resource_name"`
	RoleID     int `json:"role_id"`
	RoleName string `json:"role_name"`
}

type RoleResourcePermission struct {
	ID                      int `json:"id"`
	ResourceRoleID int `json:"resource_role_id"`
	ResourceName string  `json:"resource_name"`
	RoleName string  `json:"role_name"`
	PermissionID            int `json:"permission_id"`
	PermissionName string  `json:"permission_name"`
}




