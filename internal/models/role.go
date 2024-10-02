package models

type Role struct {
	RoleID   int    `json:"role_id"`
	RoleName string `json:"role_name"`
}

//Generate a json example with the struct Role 
// {
//  "role_id": 1,
//  "role_name": "role_name"
// }

