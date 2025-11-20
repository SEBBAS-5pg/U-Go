package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
)

// TripRepository maneja la interacción con la tabla 'trips' en PostgreSQL
type TripRepository struct {
	db *sql.DB
}

// NewTripRepository es la fábrica
func NewTripRepository(db *sql.DB) *TripRepository {
	return &TripRepository{db: db}
}

// CreateTrip registra un nuevo viaje en la base de datos
func (r *TripRepository) CreateTrip(ctx context.Context, trip *models.Trip) (*models.Trip, error) {
	// Importante: No se insertan conductor_id ni vehicle_id al inicio,
	// ya que el viaje está en estado 'solicitado'.
	query := `
		INSERT INTO trips (
			pasajero_id, 
			status, 
			origin_lat, 
			origin_lng, 
			origin_name, 
			destination_lat, 
			destination_lng, 
			destination_name
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		) RETURNING id, created_at
	`
	// Usamos un Struct para capturar el ID y CreatedAt que genera la DB
	var id string
	var createdAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		query,
		trip.PasajeroID,
		models.TripStatusSolicitado, // Estado inicial
		trip.OriginLat,
		trip.OriginLng,
		trip.OriginName,
		trip.DestinationLat,
		trip.DestinationLng,
		trip.DestinationName,
	).Scan(&id, &createdAt)

	if err != nil {
		log.Printf("Error al crear viaje en DB: %v", err)
		return nil, errors.New("error al crear el viaje en la base de datos")
	}

	// Asignar los valores devueltos al struct original
	trip.ID = id
	trip.Status = models.TripStatusSolicitado
	if createdAt.Valid {
		trip.CreatedAt = createdAt.Time
	}

	return trip, nil
}
