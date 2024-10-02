package relations

type UserRole struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

//Generate a json example with the struct UserRole
// {
//  "user_id": 1,
//  "role_id": 1
// }
