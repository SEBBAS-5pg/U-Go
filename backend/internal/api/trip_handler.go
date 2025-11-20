package api

import (
	"encoding/json"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

// TripHandler maneja las peticiones relacionadas con los viajes (creación, aceptación, etc.)
type TripHandler struct {
	tripService *service.TripService
}

// NewTripHandler es la fábrica
func NewTripHandler(tripService *service.TripService) *TripHandler {
	return &TripHandler{
		tripService: tripService,
	}
}

// CreateTrip maneja el [POST] /api/v1/trips
func (h *TripHandler) CreateTrip(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el PasajeroID del contexto (proporcionado por AuthMiddleware)
	pasajeroID, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || pasajeroID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado o ID de pasajero no encontrado"})
		return
	}

	// 2. Decodificar el JSON de la solicitud
	var req models.CreateTripRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de solicitud inválido"})
		return
	}

	// 3. Llamar al servicio para crear el viaje
	createdTrip, err := h.tripService.CreateTrip(r.Context(), &req, pasajeroID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 4. Respuesta de éxito
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTrip)
}
