package services

import (
    "database/sql"
    "policyAuth/internal/models"
)

type RoleService struct {
    DB *sql.DB
}

func NewRoleService(DB *sql.DB) *RoleService {
    return &RoleService{DB: DB}
}

func (s *RoleService) GetRoles() ([]models.Role, error) {
    rows, err := s.DB.Query("SELECT role_id, role_name FROM pds_roles")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var roles []models.Role
    for rows.Next() {
        var role models.Role
        if err := rows.Scan(&role.RoleID, &role.RoleName); err != nil {
            return nil, err
        }
        roles = append(roles, role)
    }

    return roles, nil
}

func (s *RoleService) CreateRole(role models.Role) (models.Role, error) {
    err := s.DB.QueryRow("INSERT INTO pds_roles (role_name) VALUES ($1) RETURNING role_id", role.RoleName).Scan(&role.RoleID)
    if err != nil {
        return models.Role{}, err
    }

    return role, nil
}

func (s *RoleService) UpdateRole(id int, role models.Role) error {
    _, err := s.DB.Exec("UPDATE pds_roles SET role_name=$1 WHERE role_id=$2", role.RoleName, id)
    return err
}

func (s *RoleService) DeleteRole(id int) error {
    _, err := s.DB.Exec("DELETE FROM pds_roles WHERE role_id=$1", id)
    return err
}

