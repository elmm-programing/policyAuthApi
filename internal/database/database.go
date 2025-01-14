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
    role_id INT REFERENCES pds_roles(role_id)
);

-- Table for resource-permission mappings
CREATE TABLE IF NOT EXISTS pds_resource_permission (
    id SERIAL PRIMARY KEY,
    resource_id INT REFERENCES pds_resources(resource_id),
    permission_id INT REFERENCES pds_permissions(permission_id)
);

-- New table to store permissions for each role-resource mapping
CREATE TABLE IF NOT EXISTS pds_role_resource_permissions (
    id SERIAL PRIMARY KEY,
    resource_role_id INT REFERENCES pds_resource_role(id),
    permission_id INT REFERENCES pds_permissions(permission_id)
);

INSERT INTO pds_users (username,"password")
	VALUES 
	('edwin','edwin'),
	('anonimo','anonimo'),
	('admin','admin')
ON CONFLICT (username) DO NOTHING;

INSERT INTO pds_roles (role_name)
	VALUES 
	('admin'),
	('anonimo'),
	('user'),
	('analista')
ON CONFLICT (role_name) DO NOTHING;

INSERT INTO pds_resources (resource_name)
	VALUES 
	('/home'),
	('/analistas'),
	('/normalUser'),
	('/anonimo'),
	('/ue'),
	('/clasificadores')
ON CONFLICT (resource_name) DO NOTHING;

INSERT INTO pds_permissions (permission_name)
	VALUES 
	('get'),
	('post'),
	('update'),
	('delete')
ON CONFLICT (permission_name) DO NOTHING;	

INSERT INTO pds_user_roles (user_id,role_id)
	VALUES 
	(3,1),
	(2,2),
	(1,4),
	(1,3)
ON CONFLICT (user_id,role_id) DO NOTHING;


INSERT INTO public.pds_resource_role (id,resource_id,role_id)
	VALUES (1,1,1),
	(2,1,2),
	(3,1,3),
	(4,1,4),
	(5,2,1),
	(6,2,4),
	(7,3,3),
	(8,4,2),
  (9,5,3),
  (10,5,4),
  (11,6,3),
  (12,6,4) ON CONFLICT (id) DO NOTHING;

INSERT INTO public.pds_resource_permission (id,resource_id,permission_id)
	VALUES 
	(1,1,1), 
	(2,1,2),
	(3,1,3),
	(4,1,4),
	(5,2,1), 
	(6,2,2),
	(7,2,3),
	(8,2,4),
	(9,3,1), 
	(10,3,2),
	(11,3,3),
	(12,3,4),
	(13,4,1), 
	(14,4,2),
	(15,4,3),
	(16,4,4),
  (17,5,1), 
	(18,5,2),
	(19,5,3),
	(20,5,4),
  (21,6,1), 
	(22,6,2),
	(23,6,3),
	(24,6,4)

	ON CONFLICT (id) DO NOTHING;


INSERT INTO pds_role_resource_permissions (id,resource_role_id,permission_id)
	VALUES 
	(1,1,1), 
	(2,1,2),
	(3,1,3),
	(4,1,4),
	(5,2,1), 
	(6,2,2),
	(7,2,3),
	(8,2,4),
	(9,3,1), 
	(10,3,2),
	(11,3,3),
	(12,3,4),
	(13,4,1), 
	(14,4,2),
	(15,4,3),
	(16,4,4),
	(17,5,1), 
	(18,5,2),
	(19,5,3),
	(20,5,4),
	(21,6,1), 
	(22,6,2),
	(23,6,3),
	(24,6,4),
	(25,7,1), 
	(26,7,2),
	(27,7,3),
	(28,7,4),
	(29,8,1), 
	(30,8,2),
	(31,8,3),
	(32,8,4),
  (33,9,1),
  (34,9,2),
  (35,10,3),
  (36,10,4),
  (37,11,1),
  (38,11,2),
  (39,12,3),
  (40,12,4)
	ON CONFLICT (id) DO NOTHING;



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

type DatabaseService struct {
	DB       *sql.DB
	Database string
	Password string
	Username string
	Port     string
	Host     string
	Schema   string
	Instance *sql.DB
}

func New() *DatabaseService {
	dbServices := &DatabaseService{
		Database: os.Getenv("DB_DATABASE"),
		Password: os.Getenv("DB_PASSWORD"),
		Username: os.Getenv("DB_USERNAME"),
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
		Schema:   os.Getenv("DB_SCHEMA"),
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", dbServices.Username, dbServices.Password, dbServices.Host, dbServices.Port, dbServices.Database, dbServices.Schema)
	fmt.Println("The connection string is ", connStr)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbServices.Instance = db

	return dbServices
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *DatabaseService) Health() map[string]string {
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
func (s *DatabaseService) Close() error {
	log.Printf("Disconnected from database: %s", s.Database)
	return s.DB.Close()
}
