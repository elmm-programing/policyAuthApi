package relations

type UserRole struct {
	UserID int `json:"user_id"`
	UserName string `json:"username"`
	RoleID int `json:"role_id"`
	RoleName string `json:"role_name"`
}

//Generate a json example with the struct UserRole
// {
//  "user_id": 1,
//  "role_id": 1
// }
