package api

import (
	"encoding/json"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
	"github.com/gorilla/mux"
)

type LocationHandler struct {
	repo *repository.LocationRepository
}

func NewLocationHandler(repo *repository.LocationRepository) *LocationHandler {
	return &LocationHandler{repo: repo}
}

// Estructura para recibir el JSON del frontend (body del POST)
type UpdateLocationPayload struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// Accuracy, Timestamp, etc. pueden ser ignorados por ahora si no se guardan
}

// UpdateTripLocation maneja el POST /api/v1/trip/{tripId}/location
func (h *LocationHandler) UpdateTripLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripID := vars["tripId"]

	var payload UpdateLocationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Payload inválido", http.StatusBadRequest)
		return
	}

	// Guardamos en Mongo usando la nueva función del repo
	err := h.repo.UpsertTripLocation(r.Context(), tripID, payload.Latitude, payload.Longitude)
	if err != nil {
		http.Error(w, "Error al actualizar ubicación", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// GetTripLocation maneja el GET /api/v1/trip/{tripId}/location
func (h *LocationHandler) GetTripLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripID := vars["tripId"]

	location, err := h.repo.GetTripLocation(r.Context(), tripID)
	if err != nil {
		// Si no se encuentra (mongo.ErrNoDocuments), devolvemos 404
		http.Error(w, "Ubicación no encontrada para este viaje", http.StatusNotFound)
		return
	}

	// Mapeamos a la estructura que espera el frontend (LocationReceiveService.ts)
	response := map[string]interface{}{
		"tripId":    location.TripID,
		"latitude":  location.Location.Coordinates[1], // [Lng, Lat] -> Lat es índice 1
		"longitude": location.Location.Coordinates[0], // [Lng, Lat] -> Lng es índice 0
		"timestamp": location.UpdatedAt.UnixMilli(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
