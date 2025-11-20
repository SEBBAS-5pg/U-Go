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
	tripRepo     *repository.TripRepository
	userService  *UserService
	locationRepo *repository.LocationRepository
}

// NewTripService es la fábrica
func NewTripService(tripRepo *repository.TripRepository, userService *UserService, locationRepo *repository.LocationRepository) *TripService {
	return &TripService{
		tripRepo:     tripRepo,
		userService:  userService,
		locationRepo: locationRepo,
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

// AcceptTrip maneja la lógica de aceptar el viaje y cambiar el estado del conductor
func (s *TripService) AcceptTrip(ctx context.Context, tripID string, conductorID string, vehicleID string) (*models.Trip, error) {

	// 1. Validaciones
	if vehicleID == "" {
		return nil, errors.New("debe especificar el ID del vehículo que usará")
	}

	// 2. Actualizar el viaje a 'aceptado' en PostgreSQL
	updatedTrip, err := s.tripRepo.AcceptTrip(ctx, tripID, conductorID, vehicleID)
	if err != nil {
		return nil, err
	}

	// 3. 🚗 CAMBIO CRÍTICO DE ESTADO DEL CONDUCTOR
	// El conductor pasa a 'en_viaje' en PostgreSQL
	err = s.userService.UpdateDriverStatus(ctx, updatedTrip.ConductorID, models.DriverStatusEnViaje)
	if err != nil {
		log.Printf("ADVERTENCIA CRÍTICA: No se pudo cambiar el estado del conductor %s a 'en_viaje' en PG. Error: %v", conductorID, err)
		// No revertimos el viaje, pero logueamos la advertencia.
	}

	// 4. Actualizar estado en MongoDB a 'en_viaje' (Usaremos UpsertDriverLocation
	// para que el conductor ya no aparezca como 'online' en la búsqueda de pasajeros).
	// Asumimos que la ubicación no cambia al aceptar, solo el estado.
	// Usamos las coordenadas que tenga actualmente en Mongo. Debemos implementar
	// una función para obtener su ubicación actual o actualizar solo el status.
	// **POR AHORA, SOLO ACTUALIZAMOS EL ESTADO EN PG Y DEJAMOS QUE LA BÚSQUEDA FILTRE POR STATUS.**

	// El filtro de $geoNear ya está en el LocationRepo para buscar solo 'online'.
	// Si el estado en PostgreSQL es 'en_viaje', es suficiente.

	// 5. (SIMULACIÓN DE NOTIFICACIÓN):
	log.Printf("⚠️ Simulación: Notificación enviada al pasajero %s: ¡Tu viaje fue aceptado!", updatedTrip.PasajeroID)

	return updatedTrip, nil
}

// FinalizeTrip maneja la lógica para finalizar el viaje y devolver al conductor a 'online'
func (s *TripService) FinalizeTrip(ctx context.Context, tripID string, conductorID string, req models.FinalizeTripRequest) (*models.Trip, error) {

	// 1. Actualizar el viaje a 'finalizado' en PostgreSQL
	updatedTrip, err := s.tripRepo.FinalizeTrip(ctx, tripID, conductorID)
	if err != nil {
		return nil, err
	}

	// 2. 📍 CAMBIO CRÍTICO DE ESTADO DEL CONDUCTOR EN MONGODB Y POSTGRESQL
	// El conductor regresa a 'online' y vuelve a estar disponible.

	// Actualizar estado en PostgreSQL (User.DriverStatus)
	err = s.userService.UpdateDriverStatus(ctx, updatedTrip.ConductorID, models.DriverStatusOnline)

	// Actualizar estado y ubicación en MongoDB (DriverLocation)
	// Usamos las coordenadas finales (FinalLat/FinalLng) si están disponibles,
	// o asumimos un punto para que reaparezca.
	err = s.locationRepo.UpsertDriverLocation(
		ctx,
		updatedTrip.ConductorID,
		req.FinalLat, // Latitud final
		req.FinalLng, // Longitud final
		models.DriverStatusOnline,
	)

	if err != nil {
		log.Printf("ADVERTENCIA: No se pudo devolver el conductor %s a 'online'. Error: %v", updatedTrip.ConductorID, err)
		// No detenemos el proceso si falla el estado de ubicación.
	}

	// 3. (SIMULACIÓN DE NOTIFICACIÓN):
	log.Printf("⚠️ Simulación: Notificación enviada al pasajero %s: ¡Tu viaje ha finalizado!", updatedTrip.PasajeroID)

	return updatedTrip, nil
}
