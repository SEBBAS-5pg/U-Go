// backend/internal/service/rating_service.go
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
	"github.com/google/uuid"
)

// RatingService maneja la lógica de negocio para las calificaciones
type RatingService struct {
	ratingRepo *repository.RatingRepository
	tripRepo   *repository.TripRepository // Necesitamos el tripRepo para verificar el estado del viaje
}

// NewRatingService es la fábrica
func NewRatingService(ratingRepo *repository.RatingRepository, tripRepo *repository.TripRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		tripRepo:   tripRepo,
	}
}

// CreateRating valida la calificación y la registra
func (s *RatingService) CreateRating(ctx context.Context, req models.CreateRatingRequest, raterID uuid.UUID) (*models.Rating, error) {

	tripID, err := uuid.Parse(req.TripID)
	if err != nil {
		return nil, errors.New("invalid trip id format")
	}

	// 1. Obtener el viaje y validar que exista
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("error getting trip or trip not found: %w", err)
	}

	// 2. Validar que el viaje esté FINALIZADO
	if trip.Status != models.TripStatusFinalizado {
		return nil, errors.New("only finalized trips can be rated")
	}

	// 3. Validar que el calificador sea el Pasajero asignado
	raterIDStr := raterID.String()

	if trip.PasajeroID != raterIDStr {
		return nil, errors.New("only the passenger can rate this trip")
	}
	if trip.ConductorID == "" {
		return nil, errors.New("trip has no assigned driver to rate")
	}

	// 4. Crear el modelo de Rating (Convirtiendo el Comment a puntero)
	var commentPtr *string
	if req.Comment != "" {
		commentPtr = &req.Comment
	}

	newRating := &models.Rating{
		TripID:  req.TripID,
		RaterID: raterIDStr,       // Pasajero (el que califica)
		RatedID: trip.ConductorID, // Conductor (el calificado)
		Rating:  req.Rating,
		Comment: commentPtr,
	}

	// 5. Llamar al repositorio
	createdRating, err := s.ratingRepo.CreateRating(ctx, newRating)
	if err != nil {
		// Manejo de error de restricción UNIQUE (si ya calificó)
		if createdRating == nil && errors.Is(err, sql.ErrNoRows) { // Ajustar según el error exacto de tu DB por conflicto UNIQUE
			return nil, errors.New("this trip has already been rated by the passenger")
		}
		return nil, err
	}

	// 6. (PENDIENTE HU-10): Aquí se actualizaría el promedio de calificación del conductor

	return createdRating, nil
}
