package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

// AuthMiddleware es el struct que "sostiene" el secret
type AuthMiddleware struct {
	jwtSecret string
}
// NewAuthMiddleware es la "fabrica" para nuestro middleware
func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

// -- Guardian de seguridad
// Esta es ka funcion principal del middleware
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Obtiene la cabecera "Authorization"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Error: Authorization header missing")
			writeMiddlewareError(w, http.StatusUnauthorized, "Authorization token required")
			return
		}

		// Valida el formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Error: Invalid header format. Expected 'Bearer <token>'")
			writeMiddlewareError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}
		tokenString := parts[1]

		// Valida el token (usando el secret)
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// HS256 es el algoritmo que se usa en auth_service.go
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signature method")
			}
			return []byte(m.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Error: Invalid token: %v", err)
			writeMiddlewareError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// el token se validó. extrae el ID del usuario (el "sub" claim)
		userID, ok := (*claims)["sub"].(string)
		if !ok {
			log.Println("Error: Valid token but could not extract 'sub' (UserID)")
			writeMiddlewareError(w, http.StatusInternalServerError, "Error processing token")
			return
		}
		// Inyecta el ID del usuario en el contexto de la peticion
		// para que el siguiente handler (GetMyProfiile) pueda usarlo.
		ctx := context.WithValue(r.Context(), service.ContextKeyUserID, userID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// writeMiddlewareError es una funcion de ayuda para este archivo
func writeMiddlewareError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
