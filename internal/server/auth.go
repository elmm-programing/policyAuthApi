package server

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/Nerzal/gocloak/v13"
)

var (
	keycloakClient *gocloak.GoCloak
	realm          string
	clientID       string
	clientSecret   string
)

func InitKeycloak() {
	keycloakClient = gocloak.NewClient(os.Getenv("KEYCLOAK_URL"))
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

