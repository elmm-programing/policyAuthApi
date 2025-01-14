package relations

type ResourcePermission struct {
  Id int `json:"id"`
	ResourceID   int `json:"resource_id"`
	ResourceName   string `json:"resource_name"`
	PermissionID int `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}

