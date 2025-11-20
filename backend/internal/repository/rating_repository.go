// backend/internal/repository/rating_repository.go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/SEBBAS-5pg/U-Go/backend/internal/models"
)

// RatingRepository maneja la interacción con la tabla 'ratings'
type RatingRepository struct {
	db *sql.DB
}

// NewRatingRepository es la fábrica
func NewRatingRepository(db *sql.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

// CreateRating registra una nueva calificación
func (r *RatingRepository) CreateRating(ctx context.Context, rating *models.Rating) (*models.Rating, error) {
	// Columnas: trip_id, rater_id, rated_id, rating, comment
	query := `
		INSERT INTO ratings (
			trip_id, rater_id, rated_id, rating, comment
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING id, created_at
	`

	var id string
	var createdAt sql.NullTime

	// 1. Manejo del puntero *string a sql.NullString para la base de datos
	var commentValue sql.NullString
	if rating.Comment != nil {
		commentValue = sql.NullString{String: *rating.Comment, Valid: true}
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		rating.TripID,
		rating.RaterID,
		rating.RatedID,
		rating.Rating, // Usamos el campo Rating (int) de tu struct models.Rating
		commentValue,  // Usamos el valor NullString preparado
	).Scan(&id, &createdAt)

	if err != nil {
		log.Printf("Error al crear calificación en DB: %v", err)
		return nil, errors.New("error al crear la calificación en la base de datos")
	}

	rating.ID = id
	if createdAt.Valid {
		rating.CreatedAt = createdAt.Time
	}

	return rating, nil
}
