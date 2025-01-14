package relations

import (
    "database/sql"
    "policyAuth/internal/models/relations"
)

type ResourcePermissionService struct {
    DB *sql.DB
}

func NewResourcePermissionService(DB *sql.DB) *ResourcePermissionService {
    return &ResourcePermissionService{DB: DB}
}

func (s *ResourcePermissionService) GetResourcePermissions() ([]relations.ResourcePermission, error) {
    rows, err := s.DB.Query(`
        SELECT 
            rp.id,
            rp.resource_id, 
            r.resource_name, 
            rp.permission_id, 
            p.permission_name
        FROM 
            pds_resource_permission rp
        JOIN 
            pds_resources r ON rp.resource_id = r.resource_id
        JOIN 
            pds_permissions p ON rp.permission_id = p.permission_id;
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var resourcePermissions []relations.ResourcePermission
    for rows.Next() {
        var resourcePermission relations.ResourcePermission
        if err := rows.Scan(&resourcePermission.Id,&resourcePermission.ResourceID, &resourcePermission.ResourceName, &resourcePermission.PermissionID, &resourcePermission.PermissionName); err != nil {
            return nil, err
        }
        resourcePermissions = append(resourcePermissions, resourcePermission)
    }

    return resourcePermissions, nil
}

func (s *ResourcePermissionService) CreateResourcePermission(resourcePermission relations.ResourcePermission) error {
    _, err := s.DB.Exec("INSERT INTO pds_resource_permission (resource_id, permission_id) VALUES ($1, $2)", resourcePermission.ResourceID, resourcePermission.PermissionID)
    return err
}

func (s *ResourcePermissionService) DeleteResourcePermission(id string) error {
    _, err := s.DB.Exec("DELETE FROM pds_resource_permission WHERE id=$1", id)
    return err
}

