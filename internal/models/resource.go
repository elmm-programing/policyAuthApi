package models

type Resource struct {
	ResourceID   int    `json:"resource_id"`
	ResourceName string `json:"resource_name"`
}

type ResourceWithRoleRelation struct {
	ResourceRoleId   int    `json:"resource_role_id"`
	ResourceID   int    `json:"resource_id"`
	ResourceName string `json:"resource_name"`
  Permissions []string `json:"resource_permissions"`
}


//Generate a json example with the struct Resource
// {
//  "resource_id": 1,
//  "resource_name": "resource_name"
// }
