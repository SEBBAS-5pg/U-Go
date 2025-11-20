package api

import (
	"encoding/json"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
	"github.com/gorilla/mux"
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

// AcceptTrip maneja el [POST] /api/v1/trips/{tripId}/accept (HU-11)
// Solo el conductor asignado puede aceptar.
func (h *TripHandler) AcceptTrip(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el ConductorID del contexto
	conductorID, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || conductorID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado o ID de conductor no encontrado"})
		return
	}

	// 2. Obtener el tripID de la URL
	vars := mux.Vars(r)
	tripID := vars["tripId"]
	if tripID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID del viaje no proporcionado"})
		return
	}

	// 3. Decodificar el JSON (para obtener el vehicle_id)
	var req models.AcceptTripRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de solicitud inválido"})
		return
	}

	// 4. Llamar al servicio
	updatedTrip, err := h.tripService.AcceptTrip(r.Context(), tripID, conductorID, req.VehicleID)
	if err != nil {
		// Asumiendo que el servicio retorna un error claro si el viaje ya fue aceptado
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 5. Respuesta de éxito
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTrip)
}

// FinalizeTrip maneja el [POST] /api/v1/driver/trips/{tripId}/finalize (HU-12)
func (h *TripHandler) FinalizeTrip(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el ConductorID del contexto
	conductorID, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || conductorID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado o ID de conductor no encontrado"})
		return
	}

	// 2. Obtener el tripID de la URL
	vars := mux.Vars(r)
	tripID := vars["tripId"]
	if tripID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID del viaje no proporcionado"})
		return
	}

	// 3. Decodificar el JSON (para obtener las coordenadas finales)
	var req models.FinalizeTripRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// No es un error crítico si no proporciona coordenadas finales, pero es bueno validarlo.
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de solicitud inválido"})
		return
	}

	// 4. Llamar al servicio para finalizar el viaje
	updatedTrip, err := h.tripService.FinalizeTrip(r.Context(), tripID, conductorID, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Error 400 si el conductor no es el asignado o el viaje no está activo
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 5. Respuesta de éxito
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTrip)
}
