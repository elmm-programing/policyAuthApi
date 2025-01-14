package server

import (
	"policyAuth/internal/server/handlers"
	"policyAuth/internal/server/handlers/authorization"
	relations "policyAuth/internal/server/handlers/relations"
)

func (s *Server) RegisterRoutes() {
	s.app.Get("/", s.helpers.HelloWorldHandler)
	s.app.Get("/health", JWTMiddleware(s.helpers.HealthHandler))
	s.app.Post("/auth", AuthHandler)

	userHandler := handlers.NewUserHandler(s.db.Instance)
	// User routes
	s.app.Get("/users", userHandler.GetUsers)
	s.app.Post("/users", userHandler.CreateUser)
	s.app.Put("/users", userHandler.UpdateUser)
	s.app.Delete("/users", userHandler.DeleteUser)

	roleHandler := handlers.NewRoleHandler(s.db.Instance)
	// Role routes
	s.app.Get("/roles", roleHandler.GetRoles)
	s.app.Post("/roles", roleHandler.CreateRole)
	s.app.Put("/roles", roleHandler.UpdateRole)
	s.app.Delete("/roles", roleHandler.DeleteRole)

	permissionHandler := handlers.NewPermissionHandler(s.db.Instance)
	// Permission routes
	s.app.Get("/permissions", permissionHandler.GetPermissions)
	s.app.Post("/permissions", permissionHandler.CreatePermission)
	s.app.Put("/permissions", permissionHandler.UpdatePermission)
	s.app.Delete("/permissions", permissionHandler.DeletePermission)

	resourceHandler := handlers.NewResourceHandler(s.db.Instance)
	// Resource routes
	s.app.Get("/resources", resourceHandler.GetResources)
	s.app.Post("/resources", resourceHandler.CreateResource)
	s.app.Put("/resources", resourceHandler.UpdateResource)
	s.app.Delete("/resources", resourceHandler.DeleteResource)

	userRoleHandler := relations.NewUserRoleHandler(s.db.Instance)
	// userRoleHandler := &relations.UserRoleHandler{servicesRelations.UserRoleService{DB: s.db.Instance}}
	// UserRole routes
	s.app.Get("/user_roles", userRoleHandler.GetUserRoles)
	s.app.Post("/user_roles", userRoleHandler.CreateUserRoleById)
	s.app.Delete("/user_roles", userRoleHandler.DeleteUserRoleById)

	resourceRoleHandler := relations.NewResourceRoleHandler(s.db.Instance)
	// ResourceRole routes
	s.app.Get("/resource_roles", resourceRoleHandler.GetResourceRoles)
	s.app.Post("/resource_roles", resourceRoleHandler.CreateResourceRole)
	s.app.Delete("/resource_roles/:id", resourceRoleHandler.DeleteResourceRole)

	resourcePermissionHandler := relations.NewResourcePermissionHandler(s.db.Instance)
	// RoleResourcePermission routes
	s.app.Get("/resource_permissions", resourcePermissionHandler.GetResourcePermissions)
	s.app.Post("/resource_permissions", resourcePermissionHandler.CreateResourcePermission)
	s.app.Delete("/resource_permissions/:id", resourcePermissionHandler.DeleteResourcePermission)

	roleResourcePermissionHandler := relations.NewRoleResourcePermissionHandler(s.db.Instance)
	// RoleResourcePermission routes
	s.app.Get("/role_resource_permissions", roleResourcePermissionHandler.GetRoleResourcePermissions)
	s.app.Post("/role_resource_permissions", roleResourcePermissionHandler.CreateRoleResourcePermission)
	s.app.Delete("/role_resource_permissions/:id", roleResourcePermissionHandler.DeleteRoleResourcePermission)

	resourceDetailsHandler := &authorization.ResourceDetailsHandler{DB: s.db.Instance}
	s.app.Get("/resources/:username", resourceDetailsHandler.GetRolesAndPermissionsForResource)
}
