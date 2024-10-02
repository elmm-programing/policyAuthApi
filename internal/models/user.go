package models



type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
//Generate a json example with the struct User 
// {
//  "user_id": 1,
//  "username": "username",
//  "password": "password"
// }





