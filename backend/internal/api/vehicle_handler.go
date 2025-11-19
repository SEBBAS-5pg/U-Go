package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
	"github.com/gorilla/mux"
)

// VehicleHandler maneja las peticiones HTTP relacionadas con vehículos
type VehicleHandler struct {
	vehicleService *service.VehicleService
}

// NewVehicleHandler es la "fábrica"
func NewVehicleHandler(s *service.VehicleService) *VehicleHandler {
	return &VehicleHandler{vehicleService: s}
}

// VehicleRegistrationRequest define la estructura esperada del JSON de entrada
type VehicleRegistrationRequest struct {
	Plate string `json:"plate"`
	Model string `json:"model"`
	Color string `json:"color"`
}

// RegisterVehicle maneja el [POST] /api/v1/vehicles
func (h *VehicleHandler) RegisterVehicle(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el ConductorID del contexto (inyectado por el AuthMiddleware)
	conductorID, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || conductorID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "No autorizado o ID de conductor no encontrado"})
		return
	}

	// 2. Decodificar el body de la petición
	var req VehicleRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "JSON inválido o campos faltantes"})
		return
	}

	// 3. Llamar al servicio
	vehicle, err := h.vehicleService.RegisterVehicle(r.Context(), conductorID, req.Plate, req.Model, req.Color)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 4. Devolver la respuesta de éxito
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

// UploadVehicleImage maneja el [POST] /api/v1/vehicles/{vehicleId}/image
func (h *VehicleHandler) UploadVehicleImage(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el ID del vehículo de la URL (usando Gorilla Mux vars)
	vars := mux.Vars(r)
	vehicleID := vars["vehicleId"]
	if vehicleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID de vehículo faltante en la URL"})
		return
	}

	// 2. Limitar el tamaño del archivo (ejemplo: 5MB)
	r.ParseMultipartForm(5 << 20) // 5 MB

	// 3. Obtener el archivo del formulario
	file, fileHeader, err := r.FormFile("image") // 'image' es el nombre del campo en el formulario
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener el archivo 'image' del formulario"})
		return
	}
	defer file.Close()

	// 4. Llamar al servicio
	imageURL, err := h.vehicleService.UpdateVehicleImage(r.Context(), vehicleID, file, fileHeader.Filename)
	if err != nil {
		// Logica para verificar si es error de vehicleID no encontrado
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Vehículo no encontrado"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 5. Devolver la URL
	response := map[string]string{
		"message":   "Imagen del vehículo subida exitosamente",
		"image_url": imageURL,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
