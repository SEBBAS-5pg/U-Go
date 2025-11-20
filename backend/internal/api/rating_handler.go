// backend/internal/api/rating_handler.go
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/service"
	"github.com/google/uuid"
)

// RatingHandler maneja las peticiones de calificación
type RatingHandler struct {
	ratingService *service.RatingService
}

// NewRatingHandler es la fábrica
func NewRatingHandler(ratingService *service.RatingService) *RatingHandler {
	return &RatingHandler{
		ratingService: ratingService,
	}
}

// CreateRating maneja el [POST] /api/v1/ratings
func (h *RatingHandler) CreateRating(w http.ResponseWriter, r *http.Request) {

	// 1. Obtener el RaterID (Pasajero) del contexto
	raterIDStr, ok := r.Context().Value(models.ContextUserIDKey).(string)
	if !ok || raterIDStr == "" {
		RespondWithError(w, http.StatusUnauthorized, "No autorizado o ID de calificador no encontrado")
		return
	}

	raterID, err := uuid.Parse(raterIDStr)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid Rater ID format in context")
		return
	}

	// 2. Decodificar el JSON de la solicitud
	var req models.CreateRatingRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Formato de solicitud inválido")
		return
	}

	// 3. Validación de Rating (1-5)
	if req.Rating < 1 || req.Rating > 5 {
		RespondWithError(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	// 4. Llamar al servicio
	createdRating, err := h.ratingService.CreateRating(r.Context(), req, raterID)
	if err != nil {
		// Manejo de errores de lógica de negocio (400 Bad Request)
		if strings.Contains(err.Error(), "only finalized trips") ||
			strings.Contains(err.Error(), "only the passenger can rate") ||
			strings.Contains(err.Error(), "already been rated") {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Otros errores (500 Internal Server Error)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 5. Respuesta de éxito
	RespondWithJSON(w, http.StatusCreated, createdRating)
}
