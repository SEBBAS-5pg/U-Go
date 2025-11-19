package api

import (
	"encoding/json"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

// DriverHandler maneja las peticiones relacionadas con el conductor (estado, ubicación)
type DriverHandler struct {
	userService *service.UserService // Depende del servicio de usuario para actualizar el estado/ubicación
}

// NewDriverHandler es la fábrica
func NewDriverHandler(userService *service.UserService) *DriverHandler {
	return &DriverHandler{
		userService: userService,
	}
}

// UpdateLocation maneja el [POST] /api/v1/driver/location
func (h *DriverHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el ConductorID del contexto (proporcionado por AuthMiddleware)
	userID, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado o ID de conductor no encontrado"})
		return
	}

	// 2. Decodificar el JSON de la solicitud
	var req models.UpdateLocationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de solicitud inválido"})
		return
	}

	// 3. Llamar al servicio para actualizar estado (Postgres) y ubicación (Mongo)
	err = h.userService.UpdateDriverStatusAndLocation(r.Context(), userID, req)
	if err != nil {
		// Manejo de errores específicos
		if err.Error() == "Conductor no encontrado" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 4. Respuesta de éxito
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Estado y ubicación actualizados exitosamente"})
}
