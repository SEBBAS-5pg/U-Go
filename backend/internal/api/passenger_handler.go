package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
)

// PassengerHandler maneja las peticiones relacionadas con el pasajero
type PassengerHandler struct {
	userService *service.UserService // Depende del servicio de usuario para la busqueda
}

// NewPassengerHandler es la fábrica
func NewPassengerHandler(userService *service.UserService) *PassengerHandler {
	return &PassengerHandler{
		userService: userService,
	}
}

// FindNearbyDrivers maneja el [GET] /api/v1/passenger/drivers
// Espera 'lat' y 'lng' como Query Params.
func (h *PassengerHandler) FindNearbyDrivers(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener los Query Params (lat y lng)
	query := r.URL.Query()
	latStr := query.Get("lat")
	lngStr := query.Get("lng")

	if latStr == "" || lngStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Debe proporcionar latitud ('lat') y longitud ('lng')"})
		return
	}

	// 2. Convertir a float64
	latitude, errLat := strconv.ParseFloat(latStr, 64)
	longitude, errLng := strconv.ParseFloat(lngStr, 64)

	if errLat != nil || errLng != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Latitud y Longitud deben ser números válidos"})
		return
	}

	// 3. Llamar al servicio para buscar
	drivers, err := h.userService.FindNearbyDrivers(r.Context(), latitude, longitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 4. Respuesta de éxito
	// Devolver la lista de conductores encontrados (aunque esté vacía)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(drivers)
}
