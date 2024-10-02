package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

func InitSchema(db *sql.DB) {
	createTableQuery := `
-- Table for users
CREATE TABLE IF NOT EXISTS pds_users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(50) NOT NULL
);

-- Table for roles
CREATE TABLE IF NOT EXISTS pds_roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
);

-- Table for permissions
CREATE TABLE IF NOT EXISTS pds_permissions (
    permission_id SERIAL PRIMARY KEY,
    permission_name VARCHAR(50) UNIQUE NOT NULL
);

-- Table for resources
CREATE TABLE IF NOT EXISTS pds_resources (
    resource_id SERIAL PRIMARY KEY,
    resource_name VARCHAR(50) UNIQUE NOT NULL
);

-- Table for role assignments
CREATE TABLE IF NOT EXISTS pds_user_roles (
    user_id INT REFERENCES pds_users(user_id),
    role_id INT REFERENCES pds_roles(role_id),
    PRIMARY KEY (user_id, role_id)
);

-- Table for resource-role mappings
CREATE TABLE IF NOT EXISTS pds_resource_role (
    id SERIAL PRIMARY KEY,
    resource_id INT REFERENCES pds_resources(resource_id),
    role_id INT REFERENCES pds_roles(role_id),
    PRIMARY KEY (resource_id, role_id)
);

-- Table for resource-permission mappings
CREATE TABLE IF NOT EXISTS pds_resource_permission (
    resource_id INT REFERENCES pds_resources(resource_id),
    permission_id INT REFERENCES pds_permissions(permission_id),
    PRIMARY KEY (resource_id, permission_id)
);

-- New table to store permissions for each role-resource mapping
CREATE TABLE IF NOT EXISTS pds_role_resource_permissions (
    id SERIAL PRIMARY KEY,
    resource_role_permission_id INT REFERENCES pds_role_resource_permissions(id),
    permission_id INT REFERENCES pds_permissions(permission_id)
);

  INSERT INTO pds_roles (role_name) VALUES
		('admin'),
		('user')
	ON CONFLICT (role_name) DO NOTHING;
	`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

}


// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type StructService struct {
	DB *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	DBInstance *StructService
)

func New() Service {
	// Reuse Connection
	if DBInstance != nil {
		return DBInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
  	fmt.Println("The connection string is ", connStr)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	DBInstance = &StructService{
		DB: db,
	}
	return DBInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *StructService) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.DB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.DB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *StructService) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.DB.Close()
}
