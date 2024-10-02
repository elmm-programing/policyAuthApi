package server

import (
	"policyAuth/internal/database"
	"policyAuth/internal/server/handlers"
	"policyAuth/internal/server/handlers/relations"
)

func (s *Server) RegisterRoutes() {
	s.app.Get("/", s.HelloWorldHandler)
	s.app.Get("/health", JWTMiddleware(s.healthHandler))
	s.app.Post("/auth", AuthHandler)

	userHandler := &handlers.UserHandler{DB: database.DBInstance.DB}
	// User routes
	s.app.Get("/users", userHandler.GetUsers)
	s.app.Post("/users", userHandler.CreateUser)
	s.app.Put("/users", userHandler.UpdateUser)
	s.app.Delete("/users", userHandler.DeleteUser)

	roleHandler := &handlers.RoleHandler{DB: database.DBInstance.DB}
	// Role routes
	s.app.Get("/roles", roleHandler.GetRoles)
	s.app.Post("/roles", roleHandler.CreateRole)
	s.app.Put("/roles", roleHandler.UpdateRole)
	s.app.Delete("/roles", roleHandler.DeleteRole)

	permissionHandler := &handlers.PermissionHandler{DB: database.DBInstance.DB}
	// Permission routes
	s.app.Get("/permissions", permissionHandler.GetPermissions)
	s.app.Post("/permissions", permissionHandler.CreatePermission)
	s.app.Put("/permissions", permissionHandler.UpdatePermission)
	s.app.Delete("/permissions", permissionHandler.DeletePermission)

	resourceHandler := &handlers.ResourceHandler{DB: database.DBInstance.DB}
	// Resource routes
	s.app.Get("/resources", resourceHandler.GetResources)
	s.app.Post("/resources", resourceHandler.CreateResource)
	s.app.Put("/resources", resourceHandler.UpdateResource)
	s.app.Delete("/resources", resourceHandler.DeleteResource)

	userRoleHandler := &relations.UserRoleHandler{DB: database.DBInstance.DB}
	// UserRole routes
	s.app.Get("/user_roles", userRoleHandler.GetUserRoles)
	s.app.Post("/user_roles", userRoleHandler.CreateUserRole)
	s.app.Delete("/user_roles", userRoleHandler.DeleteUserRole)

	resourceRoleHandler := &relations.ResourceRoleHandler{DB: database.DBInstance.DB}
	// ResourceRole routes
	s.app.Get("/resource_roles", resourceRoleHandler.GetResourceRoles)
	s.app.Post("/resource_roles", resourceRoleHandler.CreateResourceRole)
	s.app.Delete("/resource_roles/:id", resourceRoleHandler.DeleteResourceRole)

	roleResourcePermissionHandler := &relations.RoleResourcePermissionHandler{DB: database.DBInstance.DB}
	// RoleResourcePermission routes
	s.app.Get("/role_resource_permissions", roleResourcePermissionHandler.GetRoleResourcePermissions)
	s.app.Post("/role_resource_permissions", roleResourcePermissionHandler.CreateRoleResourcePermission)
	s.app.Delete("/role_resource_permissions/:id", roleResourcePermissionHandler.DeleteRoleResourcePermission)
}
