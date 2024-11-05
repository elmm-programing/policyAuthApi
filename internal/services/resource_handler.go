package services

import (
    "database/sql"
    "policyAuth/internal/models"
)

type ResourceService struct {
    DB *sql.DB
}

func NewResourceService(DB *sql.DB) *ResourceService {
    return &ResourceService{DB: DB}
}

func (s *ResourceService) GetResources() ([]models.Resource, error) {
    rows, err := s.DB.Query("SELECT resource_id, resource_name FROM pds_resources")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var resources []models.Resource
    for rows.Next() {
        var resource models.Resource
        if err := rows.Scan(&resource.ResourceID, &resource.ResourceName); err != nil {
            return nil, err
        }
        resources = append(resources, resource)
    }

    return resources, nil
}

func (s *ResourceService) CreateResource(resource models.Resource) (models.Resource, error) {
    err := s.DB.QueryRow("INSERT INTO pds_resources (resource_name) VALUES ($1) RETURNING resource_id", resource.ResourceName).Scan(&resource.ResourceID)
    if err != nil {
        return models.Resource{}, err
    }

    return resource, nil
}

func (s *ResourceService) UpdateResource(id int, resource models.Resource) error {
    _, err := s.DB.Exec("UPDATE pds_resources SET resource_name=$1 WHERE resource_id=$2", resource.ResourceName, id)
    return err
}

func (s *ResourceService) DeleteResource(id int) error {
    _, err := s.DB.Exec("DELETE FROM pds_resources WHERE resource_id=$1", id)
    return err
}

