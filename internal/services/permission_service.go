package services

import (
    "database/sql"
    "policyAuth/internal/models"
)

type PermissionService struct {
    DB *sql.DB
}

func NewPermissionService(DB *sql.DB) *PermissionService {
    return &PermissionService{DB: DB}
}

func (s *PermissionService) GetPermissions() ([]models.Permission, error) {
    rows, err := s.DB.Query("SELECT permission_id, permission_name FROM pds_permissions")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var permissions []models.Permission
    for rows.Next() {
        var permission models.Permission
        if err := rows.Scan(&permission.PermissionID, &permission.PermissionName); err != nil {
            return nil, err
        }
        permissions = append(permissions, permission)
    }

    return permissions, nil
}

func (s *PermissionService) CreatePermission(permission models.Permission) (models.Permission, error) {
    err := s.DB.QueryRow("INSERT INTO pds_permissions (permission_name) VALUES ($1) RETURNING permission_id", permission.PermissionName).Scan(&permission.PermissionID)
    if err != nil {
        return models.Permission{}, err
    }

    return permission, nil
}

func (s *PermissionService) UpdatePermission(id int, permission models.Permission) error {
    _, err := s.DB.Exec("UPDATE pds_permissions SET permission_name=$1 WHERE permission_id=$2", permission.PermissionName, id)
    return err
}

func (s *PermissionService) DeletePermission(id int) error {
    _, err := s.DB.Exec("DELETE FROM pds_permissions WHERE permission_id=$1", id)
    return err
}

