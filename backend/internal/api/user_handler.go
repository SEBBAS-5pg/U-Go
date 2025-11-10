package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

// UserHandler maneja las peticiones HTTP para /users
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler es la "fabrica"
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMyProfile es el handler para GET /users/me
func (h *UserHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	// Extrae el userID del contexto
	// el middleware es el que lo pone alli
	userID, ok := r.Context().Value(service.ContextKeyUserID).(string)
	if !ok {
		// Esto no pasaria si el middleware estsa bien
		log.Println("Error: /users/me reached without userID in context")
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	// Llama al servicio
	user, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	// funciona
	writeJSONResponse(w, http.StatusOK, user)
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error al escribir la respuesta JSON: %v", err)
		}
	}
}
