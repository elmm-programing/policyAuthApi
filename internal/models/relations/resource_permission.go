package relations

type ResourcePermission struct {
	ResourceID   int `json:"resource_id"`
	ResourceName   string `json:"resource_name"`
	PermissionID int `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}

