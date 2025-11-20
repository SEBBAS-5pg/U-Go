// backend/internal/service/rating_service.go
package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
	"github.com/google/uuid"
)

// RatingService maneja la lógica de negocio para las calificaciones
type RatingService struct {
	ratingRepo *repository.RatingRepository
	tripRepo   *repository.TripRepository
	userRepo   *repository.UserRepository
}

// NewRatingService es la fábrica
func NewRatingService(ratingRepo *repository.RatingRepository, tripRepo *repository.TripRepository, userRepo *repository.UserRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		tripRepo:   tripRepo,
		userRepo:   userRepo,
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
		// CORRECCIÓN 1: Manejar el error de UNIQUE *antes* de cualquier otra lógica
		// Tu repositorio devuelve "error al crear la calificación en la base de datos" si falla.
		if createdRating == nil && err.Error() == "error al crear la calificación en la base de datos" {
			// Aquí puedes intentar una detección más específica de violación de restricción UNIQUE si usas librerías como go-pg o pq.
			// Por ahora, usamos el mensaje de error general del repo:
			return nil, errors.New("this trip has already been rated by the passenger")
		}
		// Si es cualquier otro error (DB down, etc.)
		return nil, err
	}

	// CORRECCIÓN 2: La goroutine solo se ejecuta si NO hubo error (err == nil)
	// El 'createdRating' ya sabemos que NO es nil aquí.

	// 6. TAREA DE LA HU-10: Actualizar el promedio de calificación del conductor
	go func() {
		ctxUpdate := context.Background()
		errUpdate := s.userRepo.UpdateAverageRating(ctxUpdate, createdRating.RatedID)
		if errUpdate != nil {
			log.Printf("HU-10 ERROR: No se pudo actualizar el promedio de rating para el conductor %s: %v", createdRating.RatedID, errUpdate)
		}
	}()

	return createdRating, nil
}
