package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"policyAuth/internal/database"

	"github.com/Nerzal/gocloak/v13"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	keycloakClient *gocloak.GoCloak
	realm          string
	clientID       string
	clientSecret   string
)

func InitKeycloak() {
	keycloakClient = gocloak.NewClient(fmt.Sprintf("%v:%v", os.Getenv("KEYCLOAK_URL"), os.Getenv("KEYCLOAK_PORT")))
	realm = os.Getenv("KEYCLOAK_REALM")
	clientID = os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		rptResult, err := keycloakClient.RetrospectToken(context.Background(), tokenString, clientID, clientSecret, realm)
		if err != nil || !*rptResult.Active {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		_, _, err = new(jwt.Parser).ParseUnverified(tokenString, claims)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var authReq AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := keycloakClient.Login(context.Background(), clientID, clientSecret, realm, authReq.Username, authReq.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Check if the user exists in the database
	db := database.New().(*database.StructService).DB
	var userID int
	err = db.QueryRow("SELECT user_id FROM pds_users WHERE username = $1", authReq.Username).Scan(&userID)
	if err == sql.ErrNoRows {
		// User does not exist, insert the user
		err = db.QueryRow(
			"INSERT INTO pds_users (username, password) VALUES ($1, $2) RETURNING user_id",
			authReq.Username, authReq.Password,
		).Scan(&userID)
		if err != nil {
			http.Error(w, "Failed to add user to the database", http.StatusInternalServerError)
			return
		}

		// Assign default role to the new user (e.g., "user")
		var roleID int
		err = db.QueryRow("SELECT role_id FROM pds_roles WHERE role_name = $1", "user").Scan(&roleID)
		if err != nil {
			http.Error(w, "Failed to find default role", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO pds_user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
		if err != nil {
			http.Error(w, "Failed to assign role to the user", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	authResp := AuthResponse{Token: token.AccessToken}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResp)
}
