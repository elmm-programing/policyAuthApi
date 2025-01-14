package relations

import (
    "database/sql"
    "policyAuth/internal/models/relations"
)

type RoleResourcePermissionService struct {
    DB *sql.DB
}

func NewRoleResourcePermissionService(DB *sql.DB) *RoleResourcePermissionService {
    return &RoleResourcePermissionService{DB: DB}
}

func (s *RoleResourcePermissionService) GetRoleResourcePermissions() ([]relations.RoleResourcePermission, error) {
    rows, err := s.DB.Query(`
        SELECT 
            rrp.id, 
            rrp.resource_role_id, 
            r.resource_name, 
            ro.role_name, 
            rrp.permission_id, 
            p.permission_name
        FROM 
            pds_role_resource_permissions rrp
        JOIN 
            pds_resource_role rr ON rrp.resource_role_id = rr.id
        JOIN 
            pds_resources r ON rr.resource_id = r.resource_id
        JOIN 
            pds_roles ro ON rr.role_id = ro.role_id
        JOIN 
            pds_permissions p ON rrp.permission_id = p.permission_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var roleResourcePermissions []relations.RoleResourcePermission
    for rows.Next() {
        var roleResourcePermission relations.RoleResourcePermission
        if err := rows.Scan(&roleResourcePermission.ID, &roleResourcePermission.ResourceRoleID, &roleResourcePermission.ResourceName, &roleResourcePermission.RoleName, &roleResourcePermission.PermissionID, &roleResourcePermission.PermissionName); err != nil {
            return nil, err
        }
        roleResourcePermissions = append(roleResourcePermissions, roleResourcePermission)
    }

    return roleResourcePermissions, nil
}

func (s *RoleResourcePermissionService) CreateRoleResourcePermission(roleResourcePermission relations.RoleResourcePermission) error {
    _, err := s.DB.Exec("INSERT INTO pds_role_resource_permissions (resource_role_id, permission_id) VALUES ($1, $2)", roleResourcePermission.ResourceRoleID, roleResourcePermission.PermissionID)
    return err
}

func (s *RoleResourcePermissionService) DeleteRoleResourcePermission(id string) error {
    _, err := s.DB.Exec("DELETE FROM pds_role_resource_permissions WHERE id=$1", id)
    return err
}

