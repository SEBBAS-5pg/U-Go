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

type TripService struct {
	tripRepo     *repository.TripRepository
	userService  *UserService
	locationRepo *repository.LocationRepository
}

func NewTripService(tripRepo *repository.TripRepository, userService *UserService, locationRepo *repository.LocationRepository) *TripService {
	return &TripService{
		tripRepo:     tripRepo,
		userService:  userService,
		locationRepo: locationRepo,
	}
}

func (s *TripService) CreateTrip(ctx context.Context, req *models.CreateTripRequest, pasajeroID string) (*models.Trip, error) {
	// 1. NO VALIDAMOS CONDUCTOR AQUÍ (Permitimos crear viaje sin conductor)

	if req.OriginLat == 0 || req.DestinationLat == 0 {
		return nil, errors.New("debe especificar coordenadas de origen y destino")
	}

	// 2. Crear objeto
	newTrip := &models.Trip{
		PasajeroID:      pasajeroID,
		ConductorID:     "", // Vacío inicial
		OriginLat:       req.OriginLat,
		OriginLng:       req.OriginLng,
		OriginName:      req.OriginName,
		DestinationLat:  req.DestinationLat,
		DestinationLng:  req.DestinationLng,
		DestinationName: req.DestinationName,
	}

	// 3. Guardar en Repo
	createdTrip, err := s.tripRepo.CreateTrip(ctx, newTrip)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Viaje %s creado (solicitado).", createdTrip.ID)
	return createdTrip, nil
}

func (s *TripService) AcceptTrip(ctx context.Context, tripID string, conductorID string, vehicleID string) (*models.Trip, error) {
	updatedTrip, err := s.tripRepo.AcceptTrip(ctx, tripID, conductorID, vehicleID)
	if err != nil {
		return nil, err
	}
	// Actualizar estado conductor (ignorar error user service para no bloquear)
	_ = s.userService.UpdateDriverStatus(ctx, updatedTrip.ConductorID, models.DriverStatusEnViaje)
	return updatedTrip, nil
}

func (s *TripService) FinalizeTrip(ctx context.Context, tripID string, conductorID string, req models.FinalizeTripRequest) (*models.Trip, error) {
	updatedTrip, err := s.tripRepo.FinalizeTrip(ctx, tripID, conductorID)
	if err != nil {
		return nil, err
	}
	// Liberar conductor
	_ = s.userService.UpdateDriverStatus(ctx, updatedTrip.ConductorID, models.DriverStatusOnline)

	// Actualizar posición final (opcional)
	if req.FinalLat != 0 {
		_ = s.locationRepo.UpsertDriverLocation(ctx, updatedTrip.ConductorID, req.FinalLat, req.FinalLng, models.DriverStatusOnline)
	}
	return updatedTrip, nil
}

func (s *TripService) CancelTrip(ctx context.Context, tripID uuid.UUID, userID uuid.UUID) error {
	trip, err := s.tripRepo.GetByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("trip not found: %w", err)
	}
	userIDStr := userID.String()
	if trip.PasajeroID != userIDStr && (trip.ConductorID == "" || trip.ConductorID != userIDStr) {
		return errors.New("permiso denegado")
	}
	if err := s.tripRepo.CancelTrip(ctx, tripID); err != nil {
		return err
	}
	if trip.ConductorID != "" {
		_ = s.userService.UpdateDriverStatus(ctx, trip.ConductorID, "online")
	}
	return nil
}

func (s *TripService) GetPassengerHistory(ctx context.Context, pasajeroID uuid.UUID) ([]models.Trip, error) {
	return s.tripRepo.GetHistoryByPasajeroID(ctx, pasajeroID)
}

func (s *TripService) GetTripByID(ctx context.Context, tripID uuid.UUID, userID uuid.UUID) (*models.Trip, error) {
	return s.tripRepo.GetByID(ctx, tripID)
}
