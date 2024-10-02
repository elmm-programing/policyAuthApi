package models

type Permission struct {
	PermissionID   int    `json:"permission_id"`
	PermissionName string `json:"permission_name"`
}

//Generate a json example with the struct Permission
// {
//  "permission_id": 1,
//  "permission_name": "permission_name"
// }
