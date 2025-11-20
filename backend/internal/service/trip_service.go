package service

import (
	"context"
	"errors"
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/SEBBAS-5pg/U-Go/backend/internal/repository"
)

// TripService maneja la lógica de negocio para los viajes
type TripService struct {
	tripRepo *repository.TripRepository
}

// NewTripService es la fábrica
func NewTripService(tripRepo *repository.TripRepository) *TripService {
	return &TripService{
		tripRepo: tripRepo,
	}
}

// CreateTrip valida y coordina la creación de un viaje
func (s *TripService) CreateTrip(ctx context.Context, req *models.CreateTripRequest, pasajeroID string) (*models.Trip, error) {

	// 1. Validaciones de Lógica de Negocio
	if req.ConductorID == "" {
		return nil, errors.New("debe especificar el ID del conductor")
	}
	if req.OriginLat == 0 || req.DestinationLat == 0 {
		return nil, errors.New("debe especificar origen y destino")
	}

	// 2. Crear el objeto Trip para el repositorio
	newTrip := &models.Trip{
		PasajeroID:      pasajeroID,
		ConductorID:     req.ConductorID, // Se asigna el ID, aunque el conductor no haya aceptado
		OriginLat:       req.OriginLat,
		OriginLng:       req.OriginLng,
		OriginName:      req.OriginName,
		DestinationLat:  req.DestinationLat,
		DestinationLng:  req.DestinationLng,
		DestinationName: req.DestinationName,
		// Status se asigna en el repositorio como 'solicitado'
	}

	// 3. Llamar al repositorio
	createdTrip, err := s.tripRepo.CreateTrip(ctx, newTrip)
	if err != nil {
		return nil, err
	}

	// 4. (SIMULACIÓN DE NOTIFICACIÓN):
	// Aquí se integraría la lógica de WebSockets para enviar una notificación
	// al conductorID sobre el nuevo viaje.
	log.Printf("⚠️ Simulación: Notificación enviada al conductor %s para el viaje %s.", createdTrip.ConductorID, createdTrip.ID)

	return createdTrip, nil
}
