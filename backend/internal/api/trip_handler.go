package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
	"github.com/google/uuid"
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

// CancelTrip maneja POST /trips/{tripId}/cancel
func (h *TripHandler) CancelTrip(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer tripID de la URL
	vars := mux.Vars(r)
	tripIDStr := vars["tripId"]

	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		// CORRECCIÓN 1: Usar RespondWithError directamente
		RespondWithError(w, http.StatusBadRequest, "Invalid Trip ID format")
		return
	}

	// 2. Extraer userID del contexto (middleware)
	// CORRECCIÓN 2: Usar la clave de contexto definida en models (o donde sea que esté definida la clave de usuario)
	// Vemos que CreateTrip, AcceptTrip y FinalizeTrip usan models.ContextUserIDKey
	userIDStr, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || userIDStr == "" {
		// Usamos el manejo de error existente para consistencia
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado o ID de usuario no encontrado"})
		return
	}

	// Convertir el ID de usuario a UUID para el servicio
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid User ID format in context")
		return
	}

	// 3. Llamar al servicio
	if err := h.tripService.CancelTrip(r.Context(), tripID, userID); err != nil {
		// Aquí manejamos los errores específicos del servicio
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "only the assigned") {
			// CORRECCIÓN 3: Usar RespondWithError directamente
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		// CORRECCIÓN 4: Usar RespondWithError directamente
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// CORRECCIÓN 5: Usar RespondWithJSON directamente
	RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Trip successfully canceled", "status": "cancelado"})
}

// GetPassengerHistory maneja el [GET] /api/v1/passenger/history
func (h *TripHandler) GetPassengerHistory(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el UserID (PasajeroID) del contexto
	userIDStr, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || userIDStr == "" {
		RespondWithError(w, http.StatusUnauthorized, "No autorizado o ID de pasajero no encontrado")
		return
	}

	// 2. Parsear el ID a UUID
	pasajeroID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID de usuario malformado")
		return
	}

	// 3. Llamar al servicio
	history, err := h.tripService.GetPassengerHistory(r.Context(), pasajeroID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Respuesta exitosa
	RespondWithJSON(w, http.StatusOK, history)
}

// GetTripByID maneja el [GET] /api/v1/trips/{tripId}
func (h *TripHandler) GetTripByID(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el tripID de la URL
	vars := mux.Vars(r)
	tripIDStr := vars["tripId"]

	tripID, err := uuid.Parse(tripIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Formato de ID de viaje inválido")
		return
	}

	// 2. Obtener el userID del contexto (para verificación de acceso)
	userIDStr, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || userIDStr == "" {
		RespondWithError(w, http.StatusUnauthorized, "No autorizado o ID de usuario no encontrado")
		return
	}

	// Convertir el ID de usuario a UUID para el servicio
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Formato de ID de usuario inválido en contexto")
		return
	}

	// 3. Llamar al servicio
	trip, err := h.tripService.GetTripByID(r.Context(), tripID, userID)
	if err != nil {
		// El servicio maneja errores como "No encontrado" o "Acceso denegado"
		if strings.Contains(err.Error(), "not found") {
			RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}
		if strings.Contains(err.Error(), "no puede") {
			RespondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Respuesta de éxito
	RespondWithJSON(w, http.StatusOK, trip)
}
