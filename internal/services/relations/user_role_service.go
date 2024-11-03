package relations

import (
	"database/sql"

	models "policyAuth/internal/models/relations"
)

type UserRoleService struct {
	DB *sql.DB
}


func (s *UserRoleService) UserExistsById(userID int) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pds_users WHERE user_id=$1)", userID).Scan(&exists)
	return exists, err
}

func (s *UserRoleService) RoleExistsById(roleID int) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pds_roles WHERE role_id=$1)", roleID).Scan(&exists)
	return exists, err
}

func (s *UserRoleService) UserExistsByName(userName string) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pds_users WHERE username=$1)", userName).Scan(&exists)
	return exists, err
}

func (s *UserRoleService) RoleExistsByName(roleName string) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pds_roles WHERE role_name=$1)", roleName).Scan(&exists)
	return exists, err
}

func (s *UserRoleService) GetUserRoles() ([]models.UserRole, error) {
	rows, err := s.DB.Query(`
        SELECT rl.user_id, pu.username , rl.role_id, pr.role_name 
        FROM pds_user_roles rl
        join pds_roles pr on pr.role_id = rl.role_id
        join pds_users pu on pu.user_id = rl.user_id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []models.UserRole
	for rows.Next() {
		var userRole models.UserRole
		if err := rows.Scan(&userRole.UserID, &userRole.UserName, &userRole.RoleID, &userRole.RoleName); err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

func (s *UserRoleService) CreateUserRole(userRole models.UserRole) error {
	_, err := s.DB.Exec("INSERT INTO pds_user_roles (user_id, role_id) VALUES ($1, $2)", userRole.UserID, userRole.RoleID)
	return err
}

func (s *UserRoleService) DeleteUserRole(userRole models.UserRole) error {
	_, err := s.DB.Exec("DELETE FROM pds_user_roles WHERE user_id=$1 AND role_id=$2", userRole.UserID, userRole.RoleID)
	return err
}
