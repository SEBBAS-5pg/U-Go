package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
	"github.com/google/uuid"
)

// VehicleRepository maneja la interacción con la tabla 'vehicles'
type VehicleRepository struct {
	DB *sql.DB
}

// NewVehicleRepository es la fábrica
func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{DB: db}
}

// CreateVehicle inserta un nuevo vehículo en la base de datos
func (r *VehicleRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
	// Definir el ID del vehículo y el status inicial
	vehicle.ID = uuid.New().String()
	vehicle.Status = "pendiente"
	vehicle.CreatedAt = time.Now()

	query := `
		INSERT INTO vehicles (id, conductor_id, plate, model, color, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, conductor_id, plate, model, color, status, created_at
	`

	// Ejecutar la query de inserción
	err := r.DB.QueryRowContext(ctx, query,
		vehicle.ID,
		vehicle.ConductorID,
		vehicle.Plate,
		vehicle.Model,
		vehicle.Color,
		vehicle.Status,
		vehicle.CreatedAt,
	).Scan(
		&vehicle.ID,
		&vehicle.ConductorID,
		&vehicle.Plate,
		&vehicle.Model,
		&vehicle.Color,
		&vehicle.Status,
		&vehicle.CreatedAt,
	)

	if err != nil {
		log.Printf("Error al crear vehículo en DB: %v", err)
		return nil, err
	}

	return vehicle, nil
}

// GetVehiclesByConductorID recupera todos los vehículos registrados por un conductor específico
func (r *VehicleRepository) GetVehiclesByConductorID(ctx context.Context, conductorID string) ([]models.Vehicle, error) {
	vehicles := []models.Vehicle{}

	query := `
		SELECT 
			id, conductor_id, plate, model, color, status, vehicle_image_url, created_at
		FROM vehicles
		WHERE conductor_id = $1
	`

	rows, err := r.DB.QueryContext(ctx, query, conductorID)
	if err != nil {
		log.Printf("Error al buscar vehículos en DB: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v models.Vehicle
		// Escanear los resultados. Asegúrate de que el orden de las columnas coincida con el SELECT
		err := rows.Scan(
			&v.ID,
			&v.ConductorID,
			&v.Plate,
			&v.Model,
			&v.Color,
			&v.Status,
			&v.VehicleImageURL, // Puntero *string (NULLABLE)
			&v.CreatedAt,
		)
		if err != nil {
			log.Printf("Error al escanear vehículo: %v", err)
			return nil, err
		}
		vehicles = append(vehicles, v)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error después de iterar filas: %v", err)
		return nil, err
	}

	return vehicles, nil
}

// UpdateVehicleImageURL actualiza la URL/ID de la imagen del vehículo en PostgreSQL
func (r *VehicleRepository) UpdateVehicleImageURL(ctx context.Context, vehicleID string, imageURL string) error {
	query := `
		UPDATE vehicles
		SET vehicle_image_url = $1
		WHERE id = $2
	`

	result, err := r.DB.ExecContext(ctx, query, imageURL, vehicleID)
	if err != nil {
		log.Printf("Error al actualizar la URL de la imagen del vehículo en DB: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error al obtener filas afectadas después de actualizar imagen: %v", err)
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Retornar un error estándar si el ID no existe
	}

	log.Printf("✅ Imagen de vehículo actualizada para ID: %s. URL: %s", vehicleID, imageURL)
	return nil
}

// GetByID recupera un vehículo específico por su ID
func (r *VehicleRepository) GetByID(ctx context.Context, vehicleID string) (*models.Vehicle, error) {
	var v models.Vehicle

	query := `
        SELECT 
            id, conductor_id, plate, model, color, status, vehicle_image_url, created_at
        FROM vehicles
        WHERE id = $1
    `

	err := r.DB.QueryRowContext(ctx, query, vehicleID).Scan(
		&v.ID,
		&v.ConductorID,
		&v.Plate,
		&v.Model,
		&v.Color,
		&v.Status,
		&v.VehicleImageURL,
		&v.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("vehículo con ID %s no encontrado", vehicleID) // Usamos fmt.Errorf para devolver un error descriptivo
		}
		log.Printf("Error al buscar vehículo por ID %s: %v", vehicleID, err)
		return nil, err
	}

	return &v, nil
}
