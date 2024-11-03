package relations

import (
    "database/sql"
    "policyAuth/internal/models/relations"
)

type ResourceRoleService struct {
    DB *sql.DB
}

func NewResourceRoleService(DB *sql.DB) *ResourceRoleService {
    return &ResourceRoleService{DB: DB}
}

func (s *ResourceRoleService) GetResourceRoles() ([]relations.ResourceRole, error) {
    rows, err := s.DB.Query(`
        SELECT rl.id, rl.resource_id, pre.resource_name, rl.role_id, pr.role_name 
        FROM pds_resource_role rl
        JOIN pds_roles pr ON pr.role_id = rl.role_id
        JOIN pds_resources pre ON pre.resource_id = rl.resource_id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var resourceRoles []relations.ResourceRole
    for rows.Next() {
        var resourceRole relations.ResourceRole
        if err := rows.Scan(&resourceRole.ID, &resourceRole.ResourceID, &resourceRole.ResourceName, &resourceRole.RoleID, &resourceRole.RoleName); err != nil {
            return nil, err
        }
        resourceRoles = append(resourceRoles, resourceRole)
    }

    return resourceRoles, nil
}

func (s *ResourceRoleService) CreateResourceRole(resourceRole relations.ResourceRole) error {
    _, err := s.DB.Exec("INSERT INTO pds_resource_role (resource_id, role_id) VALUES ($1, $2)", resourceRole.ResourceID, resourceRole.RoleID)
    return err
}

func (s *ResourceRoleService) DeleteResourceRole(id string) error {
    _, err := s.DB.Exec("DELETE FROM pds_resource_role WHERE id=$1", id)
    return err
}

