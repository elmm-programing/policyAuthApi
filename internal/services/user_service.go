package services

import (
    "database/sql"
    "policyAuth/internal/models"
)

type UserService struct {
    DB *sql.DB
}

func NewUserService(DB *sql.DB) *UserService {
    return &UserService{DB: DB}
}

func (s *UserService) GetUsers() ([]models.User, error) {
    rows, err := s.DB.Query("SELECT user_id, username, password FROM pds_users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.UserID, &user.Username, &user.Password); err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    return users, nil
}

func (s *UserService) CreateUser(user models.User) (models.User, error) {
    err := s.DB.QueryRow("INSERT INTO pds_users (username, password) VALUES ($1, $2) RETURNING user_id", user.Username, user.Password).Scan(&user.UserID)
    if err != nil {
        return models.User{}, err
    }

    return user, nil
}

func (s *UserService) UpdateUser(id int, user models.User) error {
    _, err := s.DB.Exec("UPDATE pds_users SET username=$1, password=$2 WHERE user_id=$3", user.Username, user.Password, id)
    return err
}

func (s *UserService) DeleteUser(id int) error {
    _, err := s.DB.Exec("DELETE FROM pds_users WHERE user_id=$1", id)
    return err
}

