package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
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

// UpdateMyProfile es el handler para PUT /users/me
func (h *UserHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer el userID del contexto (puesto por el middleware)
	userID, ok := r.Context().Value(service.ContextKeyUserID).(string)
	if !ok {
		log.Println("Error: /users/me (PUT) reached without userID in context")
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		return
	}

	// 2. Decodificar el JSON del body
	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body (JSON)"})
		return
	}

	// 3. Llamar al servicio para actualizar
	updatedUser, err := h.userService.UpdateUserProfile(r.Context(), userID, &req)
	if err != nil {
		writeJSONResponse(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	// 4. Éxito
	writeJSONResponse(w, http.StatusOK, updatedUser)
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error writing JSON response: %v", err)
		}
	}
}

// UploadProfileImage es el handler para POST /users/me/image
func (h *UserHandler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer el userID del contexto
	userID, ok := r.Context().Value(service.ContextKeyUserID).(string)
	if !ok {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error interno del servidor"})
		return
	}

	// 2. Parsear la Petición (¡No es JSON! Es 'multipart/form-data')
	// Limitamos la subida a 10 MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Archivo demasiado grande (Max 10MB)"})
		return
	}

	// 3. Obtener el archivo del formulario
	// 'profile_image' es el nombre del 'key' que esperamos de Postman/Frontend
	file, handler, err := r.FormFile("profile_image")
	if err != nil {
		log.Printf("Error al obtener FormFile: %v", err)
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Campo 'profile_image' no encontrado"})
		return
	}
	defer file.Close()

	// 4. Llamar al servicio
	imageURL, err := h.userService.UpdateUserProfileImage(r.Context(), userID, file, handler.Filename)
	if err != nil {
		// (El servicio ya logueó el error)
		writeJSONResponse(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	// 5. Éxito
	writeJSONResponse(w, http.StatusOK, map[string]string{"profile_image_url": imageURL})
}
