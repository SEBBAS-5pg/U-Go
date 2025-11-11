package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

// AuthHandler es el "controlador" que maneja las peticiones HTTP de auth
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler es la "fabrica" para nuestro handler
// recibe el servicio del paso anterior
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest es el struct del JSON que esperamos recibir
// desde la app
type RegisterRequest struct {
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
}

// Register es la funcion que se conectara a la ruta POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Si el JSON es inválido, respondemos con 400 Bad Request
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validacion de entrada simple
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email, fullName and password are required"})
		return
	}

	// Llamar al "crebro" (AuthService)
	//(le pasamos el contexto de la peticion y los datos limpios)
	user, err := h.authService.Register(r.Context(), req.Email, req.FullName, req.Password)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Devolvemos un 201 Created y el objeto del usuario (sin la contraseña)
	writeJSONResponse(w, http.StatusCreated, user)
}

// Funciones para T-08 (login)

// LoginRequest es el struct del JSON que esperamos para el login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse es el struct del JSON que devolveremos
type LoginResponse struct {
	Token string `json:"token"`
}

// Login es la funcion que se conectara a la ruta POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Decodificar el JSON
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	// Validacion simple
	if req.Email == "" || req.Password == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
		return
	}

	// Llamar al "cerebro" (AuthService)
	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// el servicio devuelve  "credenciales invaliads" si falla
		writeJSONResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	// Devuelve un 200 OK y el token
	response := LoginResponse{Token: token}
	writeJSONResponse(w, http.StatusOK, response)
}

// Funcion de ayuda (helper) para escribir JSON

// WriteJSONResponse es una funcion simple para estandarizar las respuestas JSON
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// si falla el encoding, loguea el error
			log.Printf("Error writing JSON response: %v", err)
		}
	}
}
