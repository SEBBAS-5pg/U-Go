package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// AcceptTrip actualiza el estado del viaje a 'aceptado' y le asigna el vehículo
func (r *TripRepository) AcceptTrip(ctx context.Context, tripID string, conductorID string, vehicleID string) (*models.Trip, error) {

	// La consulta actualiza el estado, el conductor y el vehículo
	query := `
        UPDATE trips
        SET status = $1, conductor_id = $2, vehicle_id = $3
        WHERE id = $4 AND status = 'solicitado'
        RETURNING 
            id, pasajero_id, conductor_id, vehicle_id, status, origin_lat, origin_lng, origin_name, 
            destination_lat, destination_lng, destination_name, created_at
    `

	updatedTrip := &models.Trip{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		models.TripStatusAceptado, // $1
		conductorID,               // $2
		vehicleID,                 // $3
		tripID,                    // $4
	).Scan(
		&updatedTrip.ID,
		&updatedTrip.PasajeroID,
		&updatedTrip.ConductorID,
		&updatedTrip.VehicleID,
		&updatedTrip.Status,
		&updatedTrip.OriginLat,
		&updatedTrip.OriginLng,
		&updatedTrip.OriginName,
		&updatedTrip.DestinationLat,
		&updatedTrip.DestinationLng,
		&updatedTrip.DestinationName,
		&updatedTrip.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("viaje no encontrado o ya ha sido aceptado/cancelado")
		}
		log.Printf("Error al aceptar viaje en DB: %v", err)
		return nil, errors.New("error al actualizar el estado del viaje")
	}

	return updatedTrip, nil
}

// FinalizeTrip actualiza el estado del viaje a 'finalizado' y establece la hora de finalización.
func (r *TripRepository) FinalizeTrip(ctx context.Context, tripID string, conductorID string) (*models.Trip, error) {

	// Filtro: Busca por ID de viaje, debe estar 'aceptado' (o 'en_curso') y el conductor debe ser el que finaliza.
	filter := "id = $1 AND conductor_id = $2 AND (status = $3 OR status = $4)"

	// Los campos a actualizar (usamos NOW() para completed_at)
	update := "status = $5, completed_at = NOW()"

	// Consulta SQL para actualizar y devolver la fila
	query := fmt.Sprintf(`
        UPDATE trips 
        SET %s 
        WHERE %s 
        RETURNING 
            id, pasajero_id, conductor_id, vehicle_id, status, origin_lat, origin_lng, origin_name, 
            destination_lat, destination_lng, destination_name, created_at, started_at, completed_at
        `, update, filter)

	var trip models.Trip
	// Campos que pueden ser nulos deben ser escaneados usando punteros (*time.Time)
	var startedAt, completedAt sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		query,
		tripID,
		conductorID,
		models.TripStatusAceptado,   // Estado esperado 1
		models.TripStatusEnCurso,    // Estado esperado 2
		models.TripStatusFinalizado, // Nuevo estado
	).Scan(
		&trip.ID,
		&trip.PasajeroID,
		&trip.ConductorID,
		&trip.VehicleID,
		&trip.Status,
		&trip.OriginLat,
		&trip.OriginLng,
		&trip.OriginName,
		&trip.DestinationLat,
		&trip.DestinationLng,
		&trip.DestinationName,
		&trip.CreatedAt,
		&startedAt,
		&completedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("viaje no encontrado, ya finalizado o no eres el conductor asignado")
		}
		log.Printf("Error al finalizar viaje en DB: %v", err)
		return nil, err
	}

	// Convertir sql.NullTime a *time.Time para el modelo
	if startedAt.Valid {
		trip.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		trip.CompletedAt = &completedAt.Time
	}

	return &trip, nil
}
