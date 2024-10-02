package models

type Resource struct {
	ResourceID   int    `json:"resource_id"`
	ResourceName string `json:"resource_name"`
}

//Generate a json example with the struct Resource
// {
//  "resource_id": 1,
//  "resource_name": "resource_name"
// }
